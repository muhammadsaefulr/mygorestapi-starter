package service

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/history/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/history/response"
	responses "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"

	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/history"
	mv_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/movie_details_service"
	od_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/otakudesu_scrape"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type HistoryService struct {
	Log       *logrus.Logger
	Validate  *validator.Validate
	Repo      repository.HistoryRepo
	Anim      od_service.AnimeService
	MvDetails mv_service.MovieDetailsServiceInterface
}

func NewHistoryService(repo repository.HistoryRepo, validate *validator.Validate, od_service od_service.AnimeService, mv_service mv_service.MovieDetailsServiceInterface) HistoryService {
	return HistoryService{
		Log:       utils.Log,
		Validate:  validate,
		Repo:      repo,
		Anim:      od_service,
		MvDetails: mv_service,
	}
}

func (s *HistoryService) GetAllByUserId(c *fiber.Ctx, params *request.QueryHistory) ([]response.HistoryResponse, int64, error) {

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	user := c.Locals("user").(*model.User)

	histories, total, err := s.Repo.GetAllByUserId(c.Context(), string(user.ID.String()), params)
	if err != nil {
		return nil, 0, err
	}

	var (
		result     = make([]response.HistoryResponse, 0, len(histories))
		resultLock sync.Mutex
		semaphore  = make(chan struct{}, 4)
		wg         sync.WaitGroup
	)

	for _, h := range histories {
		history := h

		semaphore <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()

			var detail *responses.MovieDetailOnlyResponse

			if strings.Contains(history.MovieId, ("/detail")) {
				history.MovieId = path.Base(history.MovieId)
			}

			movieData, err := s.MvDetails.GetById(c, history.MovieId)
			if err == nil && movieData != nil {

				detail = &responses.MovieDetailOnlyResponse{
					MovieID:      movieData.MovieID,
					MovieType:    movieData.MovieType,
					ThumbnailURL: movieData.ThumbnailURL,
					Title:        movieData.Title,
					Rating:       movieData.Rating,
					Producer:     movieData.Producer,
					Status:       movieData.Status,
					TotalEps:     strconv.Itoa(len(movieData.Episodes)),
					Studio:       movieData.Studio,
					ReleaseDate:  movieData.ReleaseDate,
					Synopsis:     movieData.Synopsis,
					Genres:       movieData.Genres,
				}
			} else {
				animeDetail, _, _, err := s.Anim.GetAnimeDetails(history.MovieId)
				if err != nil {
					s.Log.WithError(err).Warnf("Gagal ambil detail anime untuk MovieId %s", history.MovieId)
					return
				}

				detail = &responses.MovieDetailOnlyResponse{
					MovieID:      history.MovieId,
					MovieType:    "anime",
					ThumbnailURL: animeDetail.ThumbnailURL,
					Title:        animeDetail.Title,
					Rating:       animeDetail.Rating,
					Producer:     animeDetail.Producer,
					Status:       animeDetail.Status,
					TotalEps:     animeDetail.TotalEps,
					Studio:       animeDetail.Studio,
					ReleaseDate:  animeDetail.ReleaseDate,
					Synopsis:     animeDetail.Synopsis,
					Genres: func() []string {
						titles := make([]string, len(animeDetail.Genres))
						for i, g := range animeDetail.Genres {
							titles[i] = g.Title
						}
						return titles
					}(),
				}
			}

			res := convert_types.HistoryToResponse(&history, detail)

			resultLock.Lock()
			result = append(result, res)
			resultLock.Unlock()
		}()
	}

	wg.Wait()

	return result, total, nil
}

func (s *HistoryService) GetByID(c *fiber.Ctx, id uint) (*model.History, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "History not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get history by ID failed")
	}
	return data, nil
}

func (s *HistoryService) Create(c *fiber.Ctx, req *request.CreateHistory) (*model.History, error) {

	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	IdUsr, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("Error verifying token: %s", err.Error()))
	}

	req.UserId = IdUsr

	data := convert_types.CreateHistoryToModel(req)

	if err := s.Repo.Create(c.Context(), data); err != nil {
		// s.Log.Errorf("Create error: %+v", err)

		if err == gorm.ErrDuplicatedKey {
			return nil, fiber.NewError(fiber.StatusConflict, "History already exists")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create history failed")
	}

	return data, nil
}

func (s *HistoryService) Update(c *fiber.Ctx, id uint, req *request.UpdateHistory) (*model.History, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	dataGetByID, err := s.GetByID(c, id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "History not found")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get history by ID failed")
	}

	req.ID = dataGetByID.ID
	req.UserId = dataGetByID.UserId.String()
	req.MovieEpsId = dataGetByID.MovieEpsId

	data := convert_types.UpdateHistoryToModel(req)

	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update history failed")
	}
	return s.GetByID(c, id)
}

func (s *HistoryService) Delete(c *fiber.Ctx, id uint) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "History not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete history failed")
	}
	return nil
}

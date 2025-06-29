package service

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/watchlist"
	od_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/otakudesu_scrape"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type newWatchlistService struct {
	Log          *logrus.Logger
	Validate     *validator.Validate
	Repository   repository.WatchlistRepo
	AnimeService od_service.AnimeService
}

func NewWatchlistService(repo repository.WatchlistRepo, validate *validator.Validate, animeService od_service.AnimeService) WatchlistService {
	return &newWatchlistService{
		Log:          utils.Log,
		Validate:     validate,
		Repository:   repo,
		AnimeService: animeService,
	}
}

func (s *newWatchlistService) GetAllWatchlist(c *fiber.Ctx, params *request.QueryWatchlist) ([]response.WatchlistResponse, int64, error) {

	if err := s.Validate.Struct(params); err != nil {
		log.Println("Validation error:", err)
		return nil, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	IdUsr, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return nil, 0, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("Error verifying token: %s", err.Error()))
	}

	resultFromDB, totalResults, err := s.Repository.GetAllWatchlist(c.Context(), params, IdUsr)
	if err != nil {
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to get watchlist: %s", err.Error()))
	}

	var watchlistResponses []response.WatchlistResponse

	for _, watch := range resultFromDB {
		detailMovie, _, _, errAnimeService := s.AnimeService.GetAnimeDetails(watch.MovieId)
		if errAnimeService != nil {
			log.Println("Gagal ambil detail anime:", errAnimeService.Error())
			detailMovie = model.AnimeDetail{}
		}

		watchlistResponses = append(watchlistResponses, response.WatchlistResponse{
			ID:          watch.ID,
			UserId:      IdUsr,
			MovieId:     watch.MovieId,
			AnimeDetail: detailMovie,
		})
	}

	return watchlistResponses, totalResults, nil
}

func (s *newWatchlistService) CreateWatchlist(c *fiber.Ctx, req *request.CreateWatchlist) (*model.Watchlist, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	IdUsr, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("Error verifying token: %s", err.Error()))
	}

	detailMovie, _, _, errAnimeService := s.AnimeService.GetAnimeDetails(req.MovieId)

	if errAnimeService != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to get anime detail because %s", errAnimeService.Error()))
	}

	if detailMovie.Title == "" {
		return nil, fiber.NewError(fiber.StatusNotFound, "Movie not found")
	}

	data := convert_types.CreateWatchlistToModel(&request.CreateWatchlist{
		UserId:  IdUsr,
		MovieId: req.MovieId,
	})

	if err := s.Repository.CreateWatchlist(c.Context(), data); err != nil {
		s.Log.Errorf("CreateWatchlist error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create watchlist because %s", err.Error()))
	}
	return data, nil
}

func (s *newWatchlistService) UpdateWatchlist(c *fiber.Ctx, movie_id string, req *request.UpdateWatchlist) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	data := convert_types.UpdateWatchlistToModel(req)
	data.MovieId = movie_id

	if err := s.Repository.UpdateWatchlist(c.Context(), data); err != nil {
		s.Log.Errorf("UpdateWatchlist error: %+v", err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Watchlist not found")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update watchlist")
	}

	return nil
}

func (s *newWatchlistService) DeleteWatchlist(c *fiber.Ctx, movie_id string) error {

	log.Printf("Delete watchlist: %s", movie_id)

	user := c.Locals("user").(*model.User)

	if err := s.Repository.DeleteWatchlist(c.Context(), movie_id, user.ID.String()); err != nil {
		s.Log.Errorf("DeleteWatchlist error: %+v", err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Watchlist not found")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete watchlist")
	}
	return nil
}

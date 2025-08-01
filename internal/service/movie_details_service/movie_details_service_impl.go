package service

import (
	"log"
	"path"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/movie_details"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"

	svcAnilist "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/anilist_service"
	svcMdl "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/mdl_service"
	od_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/otakudesu_scrape"
	svcTmdb "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/tmdb_service"
)

type MovieDetailsService struct {
	Log        *logrus.Logger
	Validate   *validator.Validate
	Repo       repository.MovieDetailsRepo
	AnilistSvc svcAnilist.AnilistServiceInterface
	TmdbSvc    svcTmdb.TmdbServiceInterface
	MdlSvc     svcMdl.MdlServiceInterface
	OdService  od_service.AnimeService
}

func NewMovieDetailsService(repo repository.MovieDetailsRepo, validate *validator.Validate, SvcAn svcAnilist.AnilistServiceInterface, svcTmdb svcTmdb.TmdbServiceInterface, svcMdl svcMdl.MdlServiceInterface, svcOd od_service.AnimeService) MovieDetailsServiceInterface {
	return &MovieDetailsService{
		Log:        utils.Log,
		Validate:   validate,
		Repo:       repo,
		AnilistSvc: SvcAn,
		TmdbSvc:    svcTmdb,
		MdlSvc:     svcMdl,
		OdService:  svcOd,
	}
}

func (s *MovieDetailsService) GetAll(c *fiber.Ctx, params *request.QueryMovieDetails) ([]response.MovieDetailOnlyResponse, int64, error) {
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err

	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	dataRaw, total, err := s.Repo.GetAll(c.Context(), params)

	return convert_types.MovieDetailsModelToOnlyRespArr(dataRaw), int64(total), err
}

func (s *MovieDetailsService) GetById(c *fiber.Ctx, id string) (*model.MovieDetails, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "MovieDetails not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get movie_details by ID failed")
	}

	return data, nil
}

func (s *MovieDetailsService) GetByIDPreEps(c *fiber.Ctx, id string) (*response.MovieDetailsResponse, error) {
	data, err := s.Repo.GetByIDPreEps(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "MovieDetails not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get movie_details by ID failed")
	}

	resp := convert_types.MovieDetailsModelToResp(data, nil)

	if resp.Rekomend == nil || len(*resp.Rekomend) == 0 {
		title := data.Title
		params := &request.QueryMovieDetails{
			Search:   title,
			Category: "search",
			Type:     data.MovieType,
			Page:     1,
			Limit:    1,
		}
		log.Printf("Rekomendasi: %+v", resp.Rekomend)

		var results []response.MovieDetailOnlyResponse

		switch resp.MovieDetail.MovieType {
		case "anime":
			results, _, err = s.AnilistSvc.GetAll(c, convert_types.AnilistQuery(params))
		case "movie", "tv":
			results, _, _, err = s.TmdbSvc.GetAll(c, convert_types.TmdbQuery(params))
		case "kdrama":
			results, _, _, err = s.MdlSvc.GetAll(c, convert_types.MdlQuery(params))
		default:
			return resp, nil
		}

		if err == nil && len(results) > 0 {
			first := results[0]
			var detail *response.MovieDetailOnlyResponse

			switch resp.MovieDetail.MovieType {
			case "anime":
				detail, err = s.AnilistSvc.GetMovieDetailsByID(c, first.IDSource)
			case "movie", "tv":
				detail, err = s.TmdbSvc.GetDetailByID(c, first.IDSource, resp.MovieDetail.MovieType)
			case "kdrama":
				detail, err = s.MdlSvc.GetDetailByID(c, first.IDSource)
			}

			if err == nil && detail != nil && detail.Rekomend != nil {
				finalRekom := make([]response.MovieDetailOnlyResponse, 0)
				semaphore := make(chan struct{}, 2)
				wg := sync.WaitGroup{}
				mu := sync.Mutex{}
				for _, r := range *detail.Rekomend {
					r := r
					wg.Add(1)
					semaphore <- struct{}{}

					go func() {
						defer func() {
							<-semaphore
							wg.Done()
						}()

						found := false

						movieData, _, err := s.GetAll(c, &request.QueryMovieDetails{Search: r.Title})
						if err == nil && len(movieData) > 0 {
							r.MovieID = movieData[0].MovieID
							r.PathURL = "/movie/details/" + r.MovieID
							found = true
						} else if data.MovieType == "anime" {
							odData, err := s.OdService.GetAnimeByTitle(r.Title)
							if err == nil && len(odData) > 0 {
								r.MovieID = path.Base(strings.TrimSuffix(odData[0].URL, "/"))
								r.PathURL = "/otakudesu/detail/" + r.MovieID
								found = true
							}
						}

						if found {
							mu.Lock()
							finalRekom = append(finalRekom, r)
							mu.Unlock()
						}
					}()
				}

				wg.Wait()

				if len(finalRekom) > 0 {
					resp.Rekomend = &finalRekom
				}

				log.Printf("Rekomendasi: %+v", finalRekom)
			}
		}
	}

	return resp, nil
}

func (s *MovieDetailsService) Create(c *fiber.Ctx, req *request.CreateMovieDetails) (*model.MovieDetails, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.CreateMovieDetailsToModel(req)
	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)

		if err == gorm.ErrDuplicatedKey {
			return nil, fiber.NewError(fiber.StatusBadRequest, "MovieDetails With this ID already exists")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create movie_details failed")
	}
	return data, nil
}

func (s *MovieDetailsService) Update(c *fiber.Ctx, id string, req *request.UpdateMovieDetails) (*model.MovieDetails, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	data := convert_types.UpdateMovieDetailsToModel(req)
	data.MovieID = id

	dataRespUpdt, err := s.Repo.Update(c.Context(), data)

	if err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update movie_details failed")
	}

	return dataRespUpdt, err
}

func (s *MovieDetailsService) Delete(c *fiber.Ctx, id string) error {
	if _, err := s.Repo.GetByIDPreEps(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "MovieDetails not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete movie_details failed")
	}
	return nil
}

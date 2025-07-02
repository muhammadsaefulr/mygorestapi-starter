package service

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/movie_episode"
	serviceMvDtl "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/movie_details_service"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type MovieEpisodeService struct {
	Log            *logrus.Logger
	Validate       *validator.Validate
	Repo           repository.MovieEpisodeRepo
	MovieDetailSvc serviceMvDtl.MovieDetailsServiceInterface
	S3             *utils.S3Uploader
}

func NewMovieEpisodeService(repo repository.MovieEpisodeRepo, validate *validator.Validate, s3 *utils.S3Uploader, movieDetailSvc serviceMvDtl.MovieDetailsServiceInterface) MovieEpisodeServiceInterface {
	return &MovieEpisodeService{
		Log:            utils.Log,
		Validate:       validate,
		Repo:           repo,
		MovieDetailSvc: movieDetailSvc,
		S3:             s3,
	}
}

func (s *MovieEpisodeService) GetAll(c *fiber.Ctx, params *request.QueryMovieEpisode) ([]model.MovieEpisode, int64, error) {
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	return s.Repo.GetAll(c.Context(), params)
}

func (s *MovieEpisodeService) GetByID(c *fiber.Ctx, movie_eps_id string) (*model.MovieEpisode, error) {
	data, err := s.Repo.GetByID(c.Context(), movie_eps_id)

	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "MovieEpisode not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get movie_episode by ID failed")
	}

	return data, nil
}

func (s *MovieEpisodeService) GetByMovieID(c *fiber.Ctx, movie_eps_id string, movie_id string) (*response.MovieEpisodeResponses, error) {
	data, err := s.Repo.GetByMovieID(c.Context(), movie_id)

	log.Printf("params: %+v", movie_eps_id)
	log.Printf("params: %+v", movie_id)

	movieDetail, errSvc := s.MovieDetailSvc.GetById(c, movie_id)

	if errSvc != nil {
		s.Log.Errorf("GetMovieDetailsByID error: %+v", errSvc)
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "MovieEpisode not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get movie_episode by ID failed")
	}

	results := convert_types.MovieEpisodeToResp(data, *movieDetail, movie_eps_id)

	return &results, nil
}

func (s *MovieEpisodeService) Create(c *fiber.Ctx, req *request.CreateMovieEpisodes) (*model.MovieEpisode, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user := c.Locals("user").(*model.User)

	req.SourceBy = user.Name

	data := convert_types.CreateMovieEpisodesToModel(req)
	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)

		if err == gorm.ErrDuplicatedKey {
			return nil, fiber.NewError(fiber.StatusBadRequest, "MovieEpisode With this Same Data already exists")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create movie_episode failed")
	}
	return data, nil
}

func (s *MovieEpisodeService) Update(c *fiber.Ctx, movie_eps_id string, req *request.UpdateMovieEpisodes) (*model.MovieEpisode, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.UpdateMovieEpisodesToModel(req)
	data.MovieEpsID = movie_eps_id

	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update movie_episode failed")
	}

	return s.GetByID(c, movie_eps_id)
}

func (s *MovieEpisodeService) CreateUpload(c *fiber.Ctx, req *request.CreateMovieEpisodesUpload) (*model.MovieEpisode, error) {

	user := c.Locals("user").(*model.User)

	req.SourceBy = user.Name

	if err := c.BodyParser(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid form fields")
	}

	if err := s.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	movieDtl, err := s.MovieDetailSvc.GetByIDPreEps(c, req.MovieId)
	if err != nil {
		return nil, err
	}

	file, err := req.ContentUploads.Open()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to open uploaded file")
	}
	defer file.Close()

	ext := filepath.Ext(req.ContentUploads.Filename)
	key := fmt.Sprintf("%s/episodes/%s%s", req.MovieId, req.MovieEpsID, ext)
	_, url, err := s.S3.UploadFile(movieDtl.MovieDetail.MovieType, file, key, req.ContentUploads.Header.Get("Content-Type"))
	if err != nil {
		s.Log.Errorf("S3 upload error: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to upload video")
	}

	data := &model.MovieEpisode{
		MovieEpsID: req.MovieEpsID,
		MovieId:    req.MovieId,
		Title:      req.Title,
		Resolution: req.Resolution,
		VideoURL:   url,
		SourceBy:   req.SourceBy,
	}

	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("DB create failed: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to save movie episode")
	}

	return data, nil
}

func (s *MovieEpisodeService) UpdateUpload(c *fiber.Ctx, movie_eps_id string, req *request.UpdateMovieEpisodesUpload) (*model.MovieEpisode, error) {
	panic("unimplemented")
}

func (s *MovieEpisodeService) Delete(c *fiber.Ctx, movie_eps_id string) error {
	if _, err := s.Repo.GetByID(c.Context(), movie_eps_id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "MovieEpisode not found")
	}
	if err := s.Repo.Delete(c.Context(), movie_eps_id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete movie_episode failed")
	}
	return nil
}

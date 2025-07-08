package service

import (
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
)

type MovieDetailsService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.MovieDetailsRepo
}

func NewMovieDetailsService(repo repository.MovieDetailsRepo, validate *validator.Validate) MovieDetailsServiceInterface {
	return &MovieDetailsService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *MovieDetailsService) GetAll(c *fiber.Ctx, params *request.QueryMovieDetails) ([]model.MovieDetails, int64, error) {
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err // Jika MovieID dari MovieDetailService kosong, coba cari di OdService

	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	return s.Repo.GetAll(c.Context(), params)
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
	return convert_types.MovieDetailsModelToResp(data, nil), nil
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

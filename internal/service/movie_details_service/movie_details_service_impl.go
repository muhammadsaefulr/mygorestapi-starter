package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/request"
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

func NewMovieDetailsService(repo repository.MovieDetailsRepo, validate *validator.Validate) MovieDetailsService {
	return MovieDetailsService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *MovieDetailsService) GetAll(c *fiber.Ctx, params *request.QueryMovieDetails) ([]model.MovieDetails, int64, error) {
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

func (s *MovieDetailsService) GetByID(c *fiber.Ctx, id string) (*model.MovieDetails, error) {
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
	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update movie_details failed")
	}
	return s.GetByID(c, id)
}

func (s *MovieDetailsService) Delete(c *fiber.Ctx, id string) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "MovieDetails not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete movie_details failed")
	}
	return nil
}

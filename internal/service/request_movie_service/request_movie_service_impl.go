package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_movie/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_movie/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/request_movie"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type RequestMovieService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.RequestMovieRepo
}

func NewRequestMovieService(repo repository.RequestMovieRepo, validate *validator.Validate) RequestMovieService {
	return RequestMovieService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *RequestMovieService) GetAll(c *fiber.Ctx, params *request.QueryRequestMovie) ([]response.RequestMovieResponse, int64, error) {
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	results, total, err := s.Repo.GetAll(c.Context(), params)
	if err != nil {
		return nil, 0, err
	}

	var responseResults []response.RequestMovieResponse

	for _, result := range results {
		responseResults = append(responseResults, *convert_types.ModelRequestMovieToResponse(&result))
	}

	return responseResults, total, nil
}

func (s *RequestMovieService) GetByID(c *fiber.Ctx, id uint) (*response.RequestMovieResponse, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "RequestMovie not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get request_movie by ID failed")
	}

	return convert_types.ModelRequestMovieToResponse(data), nil
}

func (s *RequestMovieService) Create(c *fiber.Ctx, req *request.CreateRequestMovie) (*model.RequestMovie, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	IdUsr, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("Error verifying token: %s", err.Error()))
	}

	req.UserIdRequest = IdUsr

	data := convert_types.CreateRequestMovieToModel(req)
	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create request_movie failed")
	}
	return data, nil
}

func (s *RequestMovieService) Update(c *fiber.Ctx, id uint, req *request.UpdateRequestMovie) (*response.RequestMovieResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.UpdateRequestMovieToModel(req)
	data.ID = id
	data.UpdatedAt = time.Now()
	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update request_movie failed")
	}
	return s.GetByID(c, id)
}

func (s *RequestMovieService) Delete(c *fiber.Ctx, id uint) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "RequestMovie not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete request_movie failed")
	}
	return nil
}

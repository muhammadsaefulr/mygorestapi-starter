package service

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/history/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/history"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type HistoryService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.HistoryRepo
}

func NewHistoryService(repo repository.HistoryRepo, validate *validator.Validate) HistoryService {
	return HistoryService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *HistoryService) GetAllByUserId(c *fiber.Ctx, params *request.QueryHistory) ([]model.History, int64, error) {
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

	return s.Repo.GetAllByUserId(c.Context(), IdUsr, params)
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

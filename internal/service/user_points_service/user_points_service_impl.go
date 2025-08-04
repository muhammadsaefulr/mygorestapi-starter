package service

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_points/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/user_points"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type UserPointsService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.UserPointsRepo
}

func NewUserPointsService(repo repository.UserPointsRepo, validate *validator.Validate) UserPointsServiceInterface {
	return &UserPointsService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *UserPointsService) GetByUserID(c *fiber.Ctx, id string) (*model.UserPoints, error) {
	data, err := s.Repo.GetByUserID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "UserPoints not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get user_points by ID failed")
	}
	return data, nil
}

func (s *UserPointsService) Update(c *fiber.Ctx, req *request.UserPoints) (*model.UserPoints, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	existing, err := s.Repo.GetByUserID(c.Context(), req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			data := convert_types.CreateUserPointsToModel(req)
			if err := s.Repo.Create(c.Context(), data); err != nil {
				s.Log.Errorf("Create error: %+v", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, "Create user_points failed")
			}
			return data, nil
		}
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get user_points")
	}

	if req.TypeUpdate == "add" {
		existing.Value += req.Value
	} else if req.TypeUpdate == "subtract" {
		if existing.Value < req.Value {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Point tidak cukup")
		}
		existing.Value -= req.Value
	}

	if err := s.Repo.Update(c.Context(), existing); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update user_points failed")
	}

	updated, err := s.Repo.GetByUserID(c.Context(), req.UserId)
	if err := s.Repo.Update(c.Context(), updated); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update user_points failed")
	}

	return updated, nil
}

func (s *UserPointsService) Delete(c *fiber.Ctx, id string) error {
	if _, err := s.Repo.GetByUserID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "UserPoints not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete user_points failed")
	}
	return nil
}

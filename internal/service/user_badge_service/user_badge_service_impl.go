package service

import (
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_badge/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/user_badge"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type UserBadgeService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.UserBadgeRepo
}

func NewUserBadgeService(repo repository.UserBadgeRepo, validate *validator.Validate) UserBadgeService {
	return UserBadgeService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *UserBadgeService) GetAll(c *fiber.Ctx, params *request.QueryUserBadge) ([]model.UserBadge, int64, error) {
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

func (s *UserBadgeService) GetByID(c *fiber.Ctx, id uint) (*model.UserBadge, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "UserBadge not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get user_badge by ID failed")
	}
	return data, nil
}

func (s *UserBadgeService) Create(c *fiber.Ctx, req *request.CreateUserBadge) (*model.UserBadge, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.CreateUserBadgeToModel(req)
	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create user_badge failed")
	}
	return data, nil
}

func (s *UserBadgeService) Update(c *fiber.Ctx, id uint, req *request.UpdateUserBadge) (*model.UserBadge, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.UpdateUserBadgeToModel(req)
	data.ID = id
	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update user_badge failed")
	}
	return s.GetByID(c, id)
}

func (s *UserBadgeService) Delete(c *fiber.Ctx, id uint) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "UserBadge not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete user_badge failed")
	}
	return nil
}

func (s *UserBadgeService) GetUserBadgeInfoByUserID(c *fiber.Ctx, user_id string) ([]model.UserBadgeInfo, error) {
	data, err := s.Repo.GetUserBadgeInfoByUserID(c.Context(), user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "UserBadge not found")
	}
	if err != nil {
		s.Log.Errorf("GetUserBadgeInfoByUserID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get user badge")
	}

	log.Printf("GetUserBadgeInfoByUserID: %+v", data)
	return data, nil
}

func (s *UserBadgeService) CreateUserBadgeInfo(c *fiber.Ctx, data *request.CreateUserBadgeInfo) error {
	if err := s.Validate.Struct(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userHandledBy := c.Locals("user").(*model.User)
	data.HandledBy = userHandledBy.ID.String()

	usrBdgeInfMdl := convert_types.CreateUserInfoBadgeToModel(data)

	err := s.Repo.CreateUserBadgeInfo(c.Context(), usrBdgeInfMdl)
	log.Printf("Err Dtail: %+v", err)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return fiber.NewError(fiber.StatusConflict, "User badge in this user already exists")
	}

	if err != nil {
		s.Log.Errorf("CreateUserBadgeInfo error: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user badge")
	}

	return nil
}

func (s *UserBadgeService) UpdateUserBadgeInfo(c *fiber.Ctx, data *request.UpdateUserBadgeInfo) error {
	if err := s.Validate.Struct(data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	handledByUser := c.Locals("user").(*model.User)
	data.HandledBy = handledByUser.ID.String()

	usrBdgeInfMdl := convert_types.UpdateUserInfoBadgeToModel(data)

	err := s.Repo.UpdateUserBadgeInfo(c.Context(), usrBdgeInfMdl)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "User badge not found")
	}

	if err != nil {
		s.Log.Errorf("UpdateUserBadgeInfo error: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user badge")
	}

	return nil
}

func (s *UserBadgeService) DeleteUserBadgeInfo(c *fiber.Ctx, user_id string) error {
	err := s.Repo.DeleteUserBadgeInfo(c.Context(), user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "User badge not found")
	}
	if err != nil {
		s.Log.Errorf("DeleteUserBadgeInfo error: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete user badge")
	}
	return nil
}

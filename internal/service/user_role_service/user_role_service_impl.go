package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_role/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/user_role"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type UserRoleService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.UserRoleRepo
}

func NewUserRoleService(repo repository.UserRoleRepo, validate *validator.Validate) UserRoleService {
	return UserRoleService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *UserRoleService) GetAll(c *fiber.Ctx, params *request.QueryUserRole) ([]model.UserRole, int64, error) {
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

func (s *UserRoleService) GetByID(c *fiber.Ctx, id uint) (*model.UserRole, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "UserRole not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get user_role by ID failed")
	}
	return data, nil
}

func (s *UserRoleService) Create(c *fiber.Ctx, req *request.CreateUserRole) (*model.UserRole, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.CreateUserRoleToModel(req)
	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)

		if err == gorm.ErrDuplicatedKey {
			return nil, fiber.NewError(fiber.StatusConflict, "User Role already exists")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create user_role failed")
	}
	return data, nil
}

func (s *UserRoleService) Update(c *fiber.Ctx, id uint, req *request.UpdateUserRole) (*model.UserRole, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.UpdateUserRoleToModel(req)
	data.ID = id

	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)

		if err == gorm.ErrDuplicatedKey {
			return nil, fiber.NewError(fiber.StatusConflict, "User Role already exists")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update user_role failed")
	}
	return s.GetByID(c, id)
}

func (s *UserRoleService) Delete(c *fiber.Ctx, id uint) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "UserRole not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete user_role failed")
	}
	return nil
}

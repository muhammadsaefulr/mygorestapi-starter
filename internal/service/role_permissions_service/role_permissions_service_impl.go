package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/dto/role_permissions/request"
	model "github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/model"
	repository "github.com/muhammadsaefulr/mygorestapi-starter/internal/repository/role_permissions"
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/shared/convert_types"
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type RolePermissionsService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	Repo     repository.RolePermissionsRepo
}

func NewRolePermissionsService(repo repository.RolePermissionsRepo, validate *validator.Validate) RolePermissionsService {
	return RolePermissionsService{
		Log:      utils.Log,
		Validate: validate,
		Repo:     repo,
	}
}

func (s *RolePermissionsService) GetAll(c *fiber.Ctx, params *request.QueryRolePermissions) ([]model.RolePermissions, int64, error) {
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

func (s *RolePermissionsService) GetByID(c *fiber.Ctx, id uint) (*model.RolePermissions, error) {
	data, err := s.Repo.GetByID(c.Context(), id)
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusNotFound, "Role Permissions not found")
	}
	if err != nil {
		s.Log.Errorf("GetByID error: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Get Role Permission by ID failed")
	}
	return data, nil
}

func (s *RolePermissionsService) Create(c *fiber.Ctx, req *request.CreateRolePermissions) (*model.RolePermissions, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.CreateRolePermissionsToModel(req)
	if err := s.Repo.Create(c.Context(), data); err != nil {
		s.Log.Errorf("Create error: %+v", err)

		if err == gorm.ErrDuplicatedKey {
			return nil, fiber.NewError(fiber.StatusConflict, "Role Permissions already exists")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, "Create Role Permission failed")
	}
	return data, nil
}

func (s *RolePermissionsService) Update(c *fiber.Ctx, id uint, req *request.UpdateRolePermissions) (*model.RolePermissions, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	data := convert_types.UpdateRolePermissionsToModel(req)
	data.ID = id
	if err := s.Repo.Update(c.Context(), data); err != nil {
		s.Log.Errorf("Update error: %+v", err)

		if err == gorm.ErrDuplicatedKey {
			return nil, fiber.NewError(fiber.StatusConflict, "Role Permissions already exists")
		}

		return nil, fiber.NewError(fiber.StatusInternalServerError, "Update Role Permission failed")
	}
	return s.GetByID(c, id)
}

func (s *RolePermissionsService) Delete(c *fiber.Ctx, id uint) error {
	if _, err := s.Repo.GetByID(c.Context(), id); err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusNotFound, "RolePermissions not found")
	}
	if err := s.Repo.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Delete Role Permission failed")
	}
	return nil
}

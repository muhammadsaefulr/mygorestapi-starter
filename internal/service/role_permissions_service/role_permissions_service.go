package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/role_permissions/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type RolePermissionsServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryRolePermissions) ([]model.RolePermissions, int64, error)
	GetByID(c *fiber.Ctx, id uint) (*model.RolePermissions, error)
	Create(c *fiber.Ctx, req *request.CreateRolePermissions) (*model.RolePermissions, error)
	Update(c *fiber.Ctx, id uint, req *request.UpdateRolePermissions) (*model.RolePermissions, error)
	Delete(c *fiber.Ctx, id uint) error
}

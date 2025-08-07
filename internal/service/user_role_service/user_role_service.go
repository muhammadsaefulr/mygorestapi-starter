package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/dto/user_role/request"
	model "github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/model"
)

type UserRoleServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryUserRole) ([]model.UserRole, int64, error)
	GetByID(c *fiber.Ctx, id uint) (*model.UserRole, error)
	Create(c *fiber.Ctx, req *request.CreateUserRole) (*model.UserRole, error)
	Update(c *fiber.Ctx, id uint, req *request.UpdateUserRole) (*model.UserRole, error)
	Delete(c *fiber.Ctx, id uint) error
}

package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_badge/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type UserBadgeServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryUserBadge) ([]model.UserBadge, int64, error)
	GetByID(c *fiber.Ctx, id uint) (*model.UserBadge, error)
	Create(c *fiber.Ctx, req *request.CreateUserBadge) (*model.UserBadge, error)
	Update(c *fiber.Ctx, id uint, req *request.UpdateUserBadge) (*model.UserBadge, error)
	Delete(c *fiber.Ctx, id uint) error

	GetUserBadgeInfoByUserID(c *fiber.Ctx, user_id string) ([]model.UserBadgeInfo, error)
	CreateUserBadgeInfo(c *fiber.Ctx, data *request.CreateUserBadgeInfo) error
	UpdateUserBadgeInfo(c *fiber.Ctx, data *request.UpdateUserBadgeInfo) error
	DeleteUserBadgeInfo(c *fiber.Ctx, user_id string) error
}

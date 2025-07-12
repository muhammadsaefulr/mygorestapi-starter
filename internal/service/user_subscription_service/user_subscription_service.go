package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_subscription/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type UserSubscriptionServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryUserSubscription) ([]model.UserSubscription, int64, error)
	GetByUserID(c *fiber.Ctx, id string) (*model.UserSubscription, error)
	Create(c *fiber.Ctx, req *request.CreateUserSubscription) (*model.UserSubscription, error)
	UpdateByUserId(c *fiber.Ctx, id string, req *request.UpdateUserSubscription) (*model.UserSubscription, error)
	DeleteByUserId(c *fiber.Ctx, id string) error
}

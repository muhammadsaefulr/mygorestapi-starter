package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/subscription_plan/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type SubscriptionPlanServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QuerySubscriptionPlan) ([]model.SubscriptionPlan, int64, error)
	GetByID(c *fiber.Ctx, id uint) (*model.SubscriptionPlan, error)
	Create(c *fiber.Ctx, req *request.CreateSubscriptionPlan) (*model.SubscriptionPlan, error)
	Update(c *fiber.Ctx, id uint, req *request.UpdateSubscriptionPlan) (*model.SubscriptionPlan, error)
	Delete(c *fiber.Ctx, id uint) error
}

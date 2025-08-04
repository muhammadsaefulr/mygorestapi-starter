package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_points/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type UserPointsServiceInterface interface {
	GetByUserID(c *fiber.Ctx, id string) (*model.UserPoints, error)
	Update(c *fiber.Ctx, req *request.UserPoints) (*model.UserPoints, error)
	Delete(c *fiber.Ctx, id string) error
}

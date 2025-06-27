package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/history/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type HistoryServiceInterface interface {
	GetAllByUserId(c *fiber.Ctx, params *request.QueryHistory) ([]model.History, int64, error)
	GetByID(c *fiber.Ctx, id uint) (*model.History, error)
	Create(c *fiber.Ctx, req *request.CreateHistory) (*model.History, error)
	Update(c *fiber.Ctx, id uint, req *request.UpdateHistory) (*model.History, error)
	Delete(c *fiber.Ctx, id uint) error
}

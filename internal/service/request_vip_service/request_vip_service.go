package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_vip/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type RequestVipServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryRequestVip) ([]model.RequestVip, int64, error)
	GetByID(c *fiber.Ctx, id uint) (*model.RequestVip, error)
	Create(c *fiber.Ctx, req *request.CreateRequestVip) (*model.RequestVip, error)
	Update(c *fiber.Ctx, id uint, req *request.UpdateRequestVip) (*model.RequestVip, error)
	Delete(c *fiber.Ctx, id uint) error
}

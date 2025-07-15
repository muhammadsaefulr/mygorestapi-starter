package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/banner_app/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type BannerAppServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryBannerApp) ([]model.BannerApp, int64, error)
	GetByID(c *fiber.Ctx, id uint) (*model.BannerApp, error)
	Create(c *fiber.Ctx, req *request.CreateBannerApp) (*model.BannerApp, error)
	Update(c *fiber.Ctx, id uint, req *request.UpdateBannerApp) (*model.BannerApp, error)
	Delete(c *fiber.Ctx, id uint) error
}

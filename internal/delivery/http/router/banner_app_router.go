package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/banner_app_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/banner_app_service"
)

func BannerAppRoutes(v1 fiber.Router, c service.BannerAppService) {
	banner_appController := controller.NewBannerAppController(c)

	group := v1.Group("/app/banner")

	group.Get("/", m.Auth("getBannerApp"), banner_appController.GetAllBannerApp)
	group.Post("/", m.Auth("createBannerApp"), banner_appController.CreateBannerApp)
	group.Get("/:id", m.Auth("getBannerApp"), banner_appController.GetBannerAppByID)
	group.Put("/:id", m.Auth("updateBannerApp"), banner_appController.UpdateBannerApp)
	group.Delete("/:id", m.Auth("deleteBannerApp"), banner_appController.DeleteBannerApp)
}

package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/anilist_controller"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/anilist_service"
)

func AnilistRoutes(v1 fiber.Router, c service.AnilistServiceInterface) {
	anilistController := controller.NewAnilistController(c)

	group := v1.Group("/anilist")

	group.Get("/", anilistController.GetAllAnilist)
	group.Get("/:id", anilistController.GetAnilistByID)
}

package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/tmdb_controller"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/tmdb_service"
)

func TmdbRoutes(v1 fiber.Router, c service.TmdbServiceInterface) {
	tmdbController := controller.NewTmdbController(c)

	group := v1.Group("/tmdb")

	group.Get("/", tmdbController.GetAllTmdb)
}

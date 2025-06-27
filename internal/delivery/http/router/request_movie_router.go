package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/request_movie_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/request_movie_service"
)

func RequestMovieRoutes(v1 fiber.Router, c service.RequestMovieService) {
	request_movieController := controller.NewRequestMovieController(c)

	group := v1.Group("/request-movie")

	group.Get("/", m.Auth(), request_movieController.GetAllRequestMovie)
	group.Post("/", m.Auth(), request_movieController.CreateRequestMovie)
	group.Get("/:id", m.Auth(), request_movieController.GetRequestMovieByID)
	group.Patch("/:id", m.Auth(), request_movieController.UpdateRequestMovie)
	group.Delete("/:id", m.Auth(), request_movieController.DeleteRequestMovie)
}

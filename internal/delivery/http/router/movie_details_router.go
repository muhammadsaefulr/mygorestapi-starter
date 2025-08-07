package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/movie_details_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/movie_details_service"
)

func MovieDetailsRoutes(v1 fiber.Router, c service.MovieDetailsServiceInterface) {
	movie_detailsController := controller.NewMovieDetailsController(c)

	group := v1.Group("/movie/details")

	group.Get("/count", movie_detailsController.GetCountAllMovieDetails)
	group.Get("/", movie_detailsController.GetAllMovieDetails)
	group.Post("/", m.Auth("createMovieDetails"), movie_detailsController.CreateMovieDetails)
	group.Get("/:id", movie_detailsController.GetMovieDetailsByID)
	group.Put("/:id", m.Auth("updateMovieDetails"), movie_detailsController.UpdateMovieDetails)
	group.Delete("/:id", m.Auth("deleteMovieDetails"), movie_detailsController.DeleteMovieDetails)
}

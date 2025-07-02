package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/movie_episode_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/movie_episode_service"
)

func MovieEpisodeRoutes(v1 fiber.Router, c service.MovieEpisodeServiceInterface) {
	movie_episodeController := controller.NewMovieEpisodeController(c)

	group := v1.Group("/movie/episodes")

	// group.Get("/", movie_episodeController.GetAllMovieEpisode)
	group.Post("/", m.Auth("addMovieEps"), movie_episodeController.CreateMovieEpisodes)
	group.Post("/upload", m.Auth("addMovieEps"), movie_episodeController.CreateUpload)
	group.Get("/:movie_id/:movie_eps_id", movie_episodeController.GetMovieEpisodeByMovieID)
	group.Get("/:id", movie_episodeController.GetMovieEpisodeByID)
	group.Put("/:id", m.Auth("updateMovieEps"), movie_episodeController.UpdateMovieEpisodes)
	group.Delete("/:id", m.Auth("deleteMovieEps"), movie_episodeController.DeleteMovieEpisode)
}

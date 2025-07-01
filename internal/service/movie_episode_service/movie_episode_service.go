package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type MovieEpisodeServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryMovieEpisode) ([]model.MovieEpisode, int64, error)
	GetByID(c *fiber.Ctx, movie_eps_id string) (*model.MovieEpisode, error)
	Create(c *fiber.Ctx, req *request.CreateMovieEpisodes) (*model.MovieEpisode, error)
	Update(c *fiber.Ctx, movie_eps_id string, req *request.UpdateMovieEpisodes) (*model.MovieEpisode, error)
	CreateUpload(c *fiber.Ctx, req *request.CreateMovieEpisodesUpload) (*model.MovieEpisode, error)
	UpdateUpload(c *fiber.Ctx, movie_eps_id string, req *request.UpdateMovieEpisodesUpload) (*model.MovieEpisode, error)
	Delete(c *fiber.Ctx, movie_eps_id string) error
}

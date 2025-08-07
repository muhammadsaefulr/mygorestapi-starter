package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type MovieEpisodeServiceInterface interface {
	GetCountAll(c *fiber.Ctx) (int64, error)
	GetAll(c *fiber.Ctx, params *request.QueryMovieEpisode) ([]model.MovieEpisode, int64, error)
	GetByID(c *fiber.Ctx, movie_eps_id string) (*model.MovieEpisode, error)
	GetMovieEpsByMovieID(c *fiber.Ctx, movie_id string, param *request.QueryMovieEpisode) ([]model.MovieEpisode, int64, error)
	GetByMovieID(c *fiber.Ctx, movie_eps_id string, movie_id string, param *request.QueryMovieEpisode) (*response.MovieEpisodeResponses, error)
	Create(c *fiber.Ctx, req *request.CreateMovieEpisodes) (*model.MovieEpisode, error)
	Update(c *fiber.Ctx, movie_eps_id string, req *request.UpdateMovieEpisodes) (*model.MovieEpisode, error)
	CreateUpload(c *fiber.Ctx, req *request.CreateMovieEpisodesUpload) (*model.MovieEpisode, error)
	UpdateUpload(c *fiber.Ctx, movie_eps_id string, req *request.UpdateMovieEpisodesUpload) (*model.MovieEpisode, error)
	Delete(c *fiber.Ctx, movie_eps_id string) error
}

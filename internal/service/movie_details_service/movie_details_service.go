package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type MovieDetailsServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryMovieDetails) ([]model.MovieDetails, int64, error)
	GetByID(c *fiber.Ctx, id string) (*model.MovieDetails, error)
	Create(c *fiber.Ctx, req *request.CreateMovieDetails) (*model.MovieDetails, error)
	Update(c *fiber.Ctx, id string, req *request.UpdateMovieDetails) (*model.MovieDetails, error)
	Delete(c *fiber.Ctx, id string) error
}

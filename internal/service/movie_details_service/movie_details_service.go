package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type MovieDetailsServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryMovieDetails) ([]response.MovieDetailOnlyResponse, int64, error)
	GetByIDPreEps(c *fiber.Ctx, id string) (*response.MovieDetailsResponse, error)
	GetById(c *fiber.Ctx, id string) (*model.MovieDetails, error)
	Create(c *fiber.Ctx, req *request.CreateMovieDetails) (*model.MovieDetails, error)
	Update(c *fiber.Ctx, id string, req *request.UpdateMovieDetails) (*model.MovieDetails, error)
	Delete(c *fiber.Ctx, id string) error
}

package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/anilist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
)

type AnilistServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryAnilist) ([]response.MovieDetailOnlyResponse, int64, error)
	GetMovieDetailsByID(c *fiber.Ctx, id string) (*response.MovieDetailOnlyResponse, error)
	GetAllGenres(c *fiber.Ctx) ([]response.GenreDetail, error)
}

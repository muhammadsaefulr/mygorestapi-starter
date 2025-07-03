package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/tmdb/request"
)

type TmdbServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryTmdb) ([]response.MovieDetailOnlyResponse, int64, error)
}

package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/discovery/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
)

type DiscoveryServiceInterface interface {
	GetDiscoverDetailByTitle(c *fiber.Ctx, mediaType, slug string) (*response.MovieDetailOnlyResponse, error)
	GetDiscoverGenres(c *fiber.Ctx, params *request.QueryDiscovery) ([]response.GenreDetail, error)
	GetDiscover(c *fiber.Ctx, params *request.QueryDiscovery) ([]response.MovieDetailOnlyResponse, int64, int64, error)
}

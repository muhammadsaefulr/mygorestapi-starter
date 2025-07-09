package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/mdl/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
)

type MdlServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryMdl) ([]response.MovieDetailOnlyResponse, int64, int64, error)
	GetDetailByID(c *fiber.Ctx, id string) (*response.MovieDetailOnlyResponse, error)
}

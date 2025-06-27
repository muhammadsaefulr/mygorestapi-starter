package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_movie/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_movie/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type RequestMovieServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryRequestMovie) ([]response.RequestMovieResponse, int64, error)
	GetByID(c *fiber.Ctx, id uint) (*response.RequestMovieResponse, error)
	Create(c *fiber.Ctx, req *request.CreateRequestMovie) (*model.RequestMovie, error)
	Update(c *fiber.Ctx, id uint, req *request.UpdateRequestMovie) (*response.RequestMovieResponse, error)
	Delete(c *fiber.Ctx, id uint) error
}

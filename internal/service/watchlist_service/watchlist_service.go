package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type WatchlistService interface {
	GetAllWatchlist(c *fiber.Ctx, params *request.QueryWatchlist) ([]response.WatchlistResponse, int64, error)
	GetWatchlistByID(c *fiber.Ctx, id uint) (*model.Watchlist, error)
	CreateWatchlist(c *fiber.Ctx, req *request.CreateWatchlist) (*model.Watchlist, error)
	UpdateWatchlist(c *fiber.Ctx, id uint, req *request.UpdateWatchlist) (*model.Watchlist, error)
	DeleteWatchlist(c *fiber.Ctx, id uint) error
}

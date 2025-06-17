package router

import (
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/watchlist_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/watchlist_service"

	"github.com/gofiber/fiber/v2"
)

func WatchlistRoutes(v1 fiber.Router, u service.WatchlistService) {
	watchlistController := controller.NewWatchlistController(u)

	group := v1.Group("/watchlists")

	group.Get("/", m.Auth(), watchlistController.GetAllWatchlist)
	group.Post("/", m.Auth(), watchlistController.CreateWatchlist)
	group.Put("/:id", m.Auth(), watchlistController.UpdateWatchlist)
	group.Delete("/:id", m.Auth(), watchlistController.DeleteWatchlist)
}

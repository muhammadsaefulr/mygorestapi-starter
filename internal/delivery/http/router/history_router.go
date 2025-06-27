package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/history_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/history_service"
)

func HistoryRoutes(v1 fiber.Router, c service.HistoryService) {
	historyController := controller.NewHistoryController(c)

	group := v1.Group("/history")

	group.Get("/", m.Auth(), historyController.GetAllHistoryByUserId)
	group.Post("/", m.Auth(), historyController.CreateHistory)
	group.Get("/:id", m.Auth(), historyController.GetHistoryByID)
	group.Put("/:id", m.Auth(), historyController.UpdateHistory)
	group.Delete("/:id", m.Auth(), historyController.DeleteHistory)
}

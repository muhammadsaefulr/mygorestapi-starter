package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/user_points_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_points_service"
)

func UserPointsRoutes(v1 fiber.Router, c service.UserPointsService) {
	user_pointsController := controller.NewUserPointsController(c)

	group := v1.Group("/users/points")

	group.Post("/", m.Auth(), user_pointsController.PostUserPoints)
	group.Get("/:id", m.Auth(), user_pointsController.GetUserPointsByID)
	group.Delete("/:id", m.Auth(), user_pointsController.DeleteUserPoints)
}

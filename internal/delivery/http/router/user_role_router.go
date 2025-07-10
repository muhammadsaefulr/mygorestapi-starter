package router

import (
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/user_role_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_role_service"
	"github.com/gofiber/fiber/v2"
)

func UserRoleRoutes(v1 fiber.Router, c service.UserRoleService) {
	user_roleController := controller.NewUserRoleController(c)

	group := v1.Group("/user_roles")

	group.Get("/", m.Auth(), user_roleController.GetAllUserRole)
	group.Post("/", m.Auth(), user_roleController.CreateUserRole)
	group.Get("/:id", m.Auth(), user_roleController.GetUserRoleByID)
	group.Patch("/:id", m.Auth(), user_roleController.UpdateUserRole)
	group.Delete("/:id", m.Auth(), user_roleController.DeleteUserRole)
}

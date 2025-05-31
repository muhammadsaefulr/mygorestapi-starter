package router

import (
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/user_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"

	system_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/system_service"
	user_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_service"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(v1 fiber.Router, u user_service.UserService, t system_service.TokenService) {
	userController := controller.NewUserController(u, t)

	user := v1.Group("/users")

	user.Get("/", m.Auth(u, "getUsers"), userController.GetUsers)
	user.Post("/", m.Auth(u, "manageUsers"), userController.CreateUser)
	user.Get("/:userId", m.Auth(u, "getUsers"), userController.GetUserByID)
	user.Patch("/:userId", m.Auth(u, "manageAcc"), userController.UpdateUser)
	user.Delete("/:userId", m.Auth(u, "manageUsers"), userController.DeleteUser)

}

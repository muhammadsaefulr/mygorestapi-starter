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

	user.Get("/count", userController.GetCountAllUser)
	user.Get("/", m.Auth("getUsers"), userController.GetUsers)
	user.Get("/session", m.Auth("getUserSession"), userController.GetUserSession)
	user.Post("/", m.Auth("manageUsers"), userController.CreateUser)
	user.Get("/:userId", m.Auth("getUsers"), userController.GetUserByID)
	user.Put("/:userId", m.Auth("manageAcc"), userController.UpdateUser)
	user.Delete("/:userId", m.Auth("manageUsers"), userController.DeleteUser)
}

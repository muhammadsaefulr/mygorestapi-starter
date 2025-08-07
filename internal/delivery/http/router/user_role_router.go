package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/mygorestapi-starter/internal/delivery/http/controller/user_role_controller"
	m "github.com/muhammadsaefulr/mygorestapi-starter/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/mygorestapi-starter/internal/service/user_role_service"
)

func UserRoleRoutes(v1 fiber.Router, c service.UserRoleService) {
	user_roleController := controller.NewUserRoleController(c)

	group := v1.Group("/user/role")

	group.Get("/", m.Auth("getUserRole"), user_roleController.GetAllUserRole)
	group.Post("/", m.Auth("createUserRole"), user_roleController.CreateUserRole)
	group.Get("/:id", m.Auth("getUserRole"), user_roleController.GetUserRoleByID)
	group.Put("/:id", m.Auth("updateUserRole"), user_roleController.UpdateUserRole)
	group.Delete("/:id", m.Auth("deleteUserRole"), user_roleController.DeleteUserRole)
}

package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/role_permissions_controller"
	// m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/role_permissions_service"
)

func RolePermissionsRoutes(v1 fiber.Router, c service.RolePermissionsService) {
	role_permissionsController := controller.NewRolePermissionsController(c)

	group := v1.Group("/user/role/permissions")

	group.Get("/", role_permissionsController.GetAllRolePermissions)
	group.Post("/", role_permissionsController.CreateRolePermissions)
	group.Get("/:id", role_permissionsController.GetRolePermissionsByID)
	group.Patch("/:id", role_permissionsController.UpdateRolePermissions)
	group.Delete("/:id", role_permissionsController.DeleteRolePermissions)
}

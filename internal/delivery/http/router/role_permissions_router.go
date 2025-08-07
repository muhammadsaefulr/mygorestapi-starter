package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/mygorestapi-starter/internal/delivery/http/controller/role_permissions_controller"
	m "github.com/muhammadsaefulr/mygorestapi-starter/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/mygorestapi-starter/internal/service/role_permissions_service"
)

func RolePermissionsRoutes(v1 fiber.Router, c service.RolePermissionsService) {
	role_permissionsController := controller.NewRolePermissionsController(c)

	group := v1.Group("/user/roles/permissions")

	group.Get("/", m.Auth("getRolePermissions"), role_permissionsController.GetAllRolePermissions)
	group.Post("/", m.Auth("createRolePermissions"), role_permissionsController.CreateRolePermissions)
	group.Get("/:id", m.Auth("getRolePermissions"), role_permissionsController.GetRolePermissionsByID)
	group.Put("/:id", m.Auth("updateRolePermissions"), role_permissionsController.UpdateRolePermissions)
	group.Delete("/:id", m.Auth("deleteRolePermissions"), role_permissionsController.DeleteRolePermissions)
}

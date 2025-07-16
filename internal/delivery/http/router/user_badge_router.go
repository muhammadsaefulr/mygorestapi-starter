package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/user_badge_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_badge_service"
)

func UserBadgeRoutes(v1 fiber.Router, c service.UserBadgeService) {
	user_badgeController := controller.NewUserBadgeController(c)

	group := v1.Group("/badge")

	group.Get("/", m.Auth("getBadge"), user_badgeController.GetAllUserBadge)
	group.Post("/", m.Auth("addBadge"), user_badgeController.CreateUserBadge)
	group.Get("/:id", m.Auth("getBadge"), user_badgeController.GetUserBadgeByID)
	group.Put("/:id", m.Auth("updateBadge"), user_badgeController.UpdateUserBadge)
	group.Delete("/:id", m.Auth("deleteBadge"), user_badgeController.DeleteUserBadge)

	groupInfo := v1.Group("/user/badge")

	groupInfo.Get("/:user_id", m.Auth("getUserBadge"), user_badgeController.GetUserBadgeInfoByUserID)
	groupInfo.Post("/", m.Auth("addUserBadge"), user_badgeController.CreateUserBadgeInfo)
	groupInfo.Put("/:user_id", m.Auth("updateUserBadge"), user_badgeController.UpdateUserBadgeInfo)
	groupInfo.Delete("/:user_id", m.Auth("deleteUserBadge"), user_badgeController.DeleteUserBadgeInfo)
}

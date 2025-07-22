package router

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/request_vip_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/request_vip_service"
)

func RequestVipRoutes(v1 fiber.Router, c service.RequestVipService) {
	request_vipController := controller.NewRequestVipController(c)

	group := v1.Group("/request-vip")

	group.Get("/", m.Auth("getRequestVIP"), request_vipController.GetAllRequestVip)
	group.Post("/", m.Auth(), request_vipController.CreateRequestVip)
	group.Get("/:id", m.Auth("getRequestVIP"), request_vipController.GetRequestVipByID)
	group.Put("/:id", m.Auth("updateRequestVIP"), request_vipController.UpdateRequestVip)
	group.Delete("/:id", m.Auth("deleteRequestVIP"), request_vipController.DeleteRequestVip)
}

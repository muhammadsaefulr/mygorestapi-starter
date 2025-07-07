package router

import (
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/mdl_controller"
	// m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	"github.com/gofiber/fiber/v2"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/mdl_service"
)

func MdlRoutes(v1 fiber.Router, c service.MdlServiceInterface) {
	mdlController := controller.NewMdlController(c)

	group := v1.Group("/mdl")

	group.Get("/", mdlController.GetAllMdl)
}

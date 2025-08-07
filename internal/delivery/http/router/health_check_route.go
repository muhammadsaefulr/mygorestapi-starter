package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/delivery/http/controller"
	service "github.com/muhammadsaefulr/mygorestapi-starter/internal/service/system_service"
)

func HealthCheckRoutes(v1 fiber.Router, h service.HealthCheckService) {
	healthCheckController := controller.NewHealthCheckController(h)

	healthCheck := v1.Group("/health-check")
	healthCheck.Get("/", healthCheckController.Check)
}

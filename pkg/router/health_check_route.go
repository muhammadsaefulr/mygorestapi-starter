package router

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/pkg/delivery/http/controller"
	"github.com/muhammadsaefulr/NimeStreamAPI/pkg/service"

	"github.com/gofiber/fiber/v2"
)

func HealthCheckRoutes(v1 fiber.Router, h service.HealthCheckService) {
	healthCheckController := controller.NewHealthCheckController(h)

	healthCheck := v1.Group("/health-check")
	healthCheck.Get("/", healthCheckController.Check)
}

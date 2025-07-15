package router

import (
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/discovery_controller"
	// m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	"github.com/gofiber/fiber/v2"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/discovery_service"
)

func DiscoveryRoutes(v1 fiber.Router, c service.DiscoveryServiceInterface) {
	discoveryController := controller.NewDiscoveryController(c)

	group := v1.Group("/discovery")

	group.Get("/", discoveryController.GetDiscover)
	group.Get("/detail/:mediaType/:title", discoveryController.GetDiscoverDetailByTitle)
	group.Get("/genres", discoveryController.GetDiscoverGenres)
}

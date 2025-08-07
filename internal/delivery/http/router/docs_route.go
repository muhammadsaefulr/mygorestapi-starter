package router

import (
	// initialize the Swagger documentation
	_ "github.com/muhammadsaefulr/mygorestapi-starter/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func DocsRoutes(v1 fiber.Router) {
	docs := v1.Group("/docs")

	docs.Get("/*", swagger.HandlerDefault)
}

package router

import (
	controller "github.com/your/module/internal/delivery/http/controller/{{.Name}}_controller"
	m "github.com/your/module/internal/delivery/middleware"

	"github.com/gofiber/fiber/v2"
)

func {{.PascalName}}Routes(v1 fiber.Router) {
	{{.Name}}Controller := controller.New{{.PascalName}}Controller()

	group := v1.Group("/{{.Name}}s")

	group.Get("/", m.Auth(nil, "get{{.PascalName}}s"), {{.Name}}Controller.GetAll)
	group.Post("/", m.Auth(nil, "manage{{.PascalName}}s"), {{.Name}}Controller.Create)
	group.Get("/:id", m.Auth(nil, "get{{.PascalName}}s"), {{.Name}}Controller.GetByID)
	group.Patch("/:id", m.Auth(nil, "manage{{.PascalName}}s"), {{.Name}}Controller.Update)
	group.Delete("/:id", m.Auth(nil, "manage{{.PascalName}}s"), {{.Name}}Controller.Delete)
}

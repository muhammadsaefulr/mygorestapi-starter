package router

import (
	controller "{{ .ModulePath }}/internal/delivery/http/controller/{{.Name}}_controller"
	m "{{ .ModulePath }}/internal/delivery/middleware"
	service "{{ .ModulePath }}/internal/service/{{.Name}}_service"
	"github.com/gofiber/fiber/v2"
)

func {{.PascalName}}Routes(v1 fiber.Router, c service.{{.PascalName}}Service) {
	{{.Name}}Controller := controller.New{{.PascalName}}Controller(c)

	group := v1.Group("/{{.Name}}s")

	group.Get("/", m.Auth(), {{.Name}}Controller.GetAll{{.PascalName}})
	group.Post("/", m.Auth(), {{.Name}}Controller.Create{{.PascalName}})
	group.Get("/:id", m.Auth(), {{.Name}}Controller.Get{{.PascalName}}ByID)
	group.Patch("/:id", m.Auth(), {{.Name}}Controller.Update{{.PascalName}})
	group.Delete("/:id", m.Auth(), {{.Name}}Controller.Delete{{.PascalName}})
}

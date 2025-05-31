package handler

import (
	"github.com/gofiber/fiber/v2"
)

type {{.PascalName}}Handler struct {}

func New{{.PascalName}}Handler() *{{.PascalName}}Handler {
	return &{{.PascalName}}Handler{}
}

func (h *{{.PascalName}}Handler) GetAll(c *fiber.Ctx) error {
	return c.JSON("list of {{.Name}}")
}

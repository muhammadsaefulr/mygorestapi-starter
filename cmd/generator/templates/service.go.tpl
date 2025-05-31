package service

import (
	"github.com/gofiber/fiber/v2"
	"{{ .ModulePath }}/internal/domain/dto/{{.Name}}/request"
	"{{ .ModulePath }}/internal/domain/model/{{.Name}}"
)

type {{.PascalName}}Service interface {
	GetAll{{.PascalName}}(c *fiber.Ctx, params *request.Query{{.PascalName}}) ([]model.{{.PascalName}}, int64, error)
	Get{{.PascalName}}ByID(c *fiber.Ctx, id string) (*model.{{.PascalName}}, error)
	Create{{.PascalName}}(c *fiber.Ctx, req *request.Create{{.PascalName}}) (*model.{{.PascalName}}, error)
	Update{{.PascalName}}(c *fiber.Ctx, id string, req *request.Update{{.PascalName}}) (*model.{{.PascalName}}, error)
	Delete{{.PascalName}}(c *fiber.Ctx, id string) error
}

package service

import (
	"github.com/gofiber/fiber/v2"

	"{{.ModulePath}}/internal/domain/dto/{{.Name}}/request"
	model "{{.ModulePath}}/internal/domain/model"
)

type {{.PascalName}}ServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.Query{{.PascalName}}) ([]model.{{.PascalName}}, int64, error)
	GetByID(c *fiber.Ctx, id uint) (*model.{{.PascalName}}, error)
	Create(c *fiber.Ctx, req *request.Create{{.PascalName}}) (*model.{{.PascalName}}, error)
	Update(c *fiber.Ctx, id uint, req *request.Update{{.PascalName}}) (*model.{{.PascalName}}, error)
	Delete(c *fiber.Ctx, id uint) error
}

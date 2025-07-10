package service

import (
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/report_error/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type ReportErrorServiceInterface interface {
	GetAll(c *fiber.Ctx, params *request.QueryReportError) ([]model.ReportError, int64, error)
	GetByID(c *fiber.Ctx, id string) (*model.ReportError, error)
	Create(c *fiber.Ctx, req *request.CreateReportError) (*model.ReportError, error)
	Update(c *fiber.Ctx, id string, req *request.UpdateReportError) (*model.ReportError, error)
	Delete(c *fiber.Ctx, id string) error
}

package controller

import (
	"math"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/report_error/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/report_error_service"
)

type ReportErrorController struct {
	Service service.ReportErrorServiceInterface
}

func NewReportErrorController(service service.ReportErrorServiceInterface) *ReportErrorController {
	return &ReportErrorController{Service: service}
}

// @Tags         report_error
// @Summary      Get all report_error
// @Security     BearerAuth
// @Produce      json
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @param        type query       string false "Type Status"
// @Router       /report-episode [get]
func (h *ReportErrorController) GetAllReportError(c *fiber.Ctx) error {
	query := &request.QueryReportError{
		Page:   c.QueryInt("page", 1),
		Search: c.Query("search", ""),
		Type:   c.Query("type", ""),
		Limit:  c.QueryInt("limit", 10),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.ReportError]{
		Code:         fiber.StatusOK,
		Status:       "success",
		Message:      "Successfully retrieved data",
		Results:      data,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   int64(math.Ceil(float64(total) / float64(query.Limit))),
		TotalResults: total,
	})
}

// @Tags         report_error
// @Summary      Get a report_error by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "ReportError ID (uint)"
// @Router       /report-episode/{id} [get]
func (h *ReportErrorController) GetReportErrorByID(c *fiber.Ctx) error {
	idStr := c.Params("id")

	result, err := h.Service.GetByID(c, idStr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.ReportError]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         report_error
// @Summary      Create a new report_error
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  request.CreateReportError  true  "Request body"
// @Router       /report-episode [post]
func (h *ReportErrorController) CreateReportError(c *fiber.Ctx) error {
	req := new(request.CreateReportError)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.ReportError]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Successfuly Reported Error",
		Data:    *result,
	})
}

// @Tags         report_error
// @Summary      Update a report_error
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "ReportError ID (uint)"
// @Param        request  body  request.UpdateReportError  true  "Request body"
// @Router       /report-episode/{id} [put]
func (h *ReportErrorController) UpdateReportError(c *fiber.Ctx) error {
	idStr := c.Params("id")

	req := new(request.UpdateReportError)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, idStr, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.ReportError]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         report_error
// @Summary      Delete a report_error
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "ReportError ID (uint)"
// @Router       /report-episode/{id} [delete]
func (h *ReportErrorController) DeleteReportError(c *fiber.Ctx) error {
	idStr := c.Params("id")

	if err := h.Service.Delete(c, idStr); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data deleted successfully",
	})
}

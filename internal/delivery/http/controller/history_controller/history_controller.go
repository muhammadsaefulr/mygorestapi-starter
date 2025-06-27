package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/history/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/history_service"
)

type HistoryController struct {
	Service service.HistoryService
}

func NewHistoryController(service service.HistoryService) *HistoryController {
	return &HistoryController{Service: service}
}

// @Tags         history
// @Summary      Get all history
// @Produce      json
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Security     BearerAuth
// @Param        search query     string  false  "Search term"
// @Router       /history [get]
func (h *HistoryController) GetAllHistoryByUserId(c *fiber.Ctx) error {
	query := &request.QueryHistory{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 10),
	}

	data, total, err := h.Service.GetAllByUserId(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.History]{
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

// @Tags         history
// @Summary      Get a history by ID
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "History ID (uint)"
// @Router       /history/{id} [get]
func (h *HistoryController) GetHistoryByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	result, err := h.Service.GetByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.History]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         history
// @Summary      Create a new history
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  request.CreateHistory  true  "Request body"
// @Router       /history [post]
func (h *HistoryController) CreateHistory(c *fiber.Ctx) error {
	req := new(request.CreateHistory)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.History]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         history
// @Summary      Update a history
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  int  true  "History ID (uint)"
// @Param        request  body  request.UpdateHistory  true  "Request body"
// @Router       /history/{id} [put]
func (h *HistoryController) UpdateHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	req := new(request.UpdateHistory)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.History]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         history
// @Summary      Delete a history
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "History ID (uint)"
// @Router       /history/{id} [delete]
func (h *HistoryController) DeleteHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	if err := h.Service.Delete(c, id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data deleted successfully",
	})
}

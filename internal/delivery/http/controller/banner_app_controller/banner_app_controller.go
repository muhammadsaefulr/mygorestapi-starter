package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/banner_app/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/banner_app_service"
)

type BannerAppController struct {
	Service service.BannerAppService
}

func NewBannerAppController(service service.BannerAppService) *BannerAppController {
	return &BannerAppController{Service: service}
}

// @Tags         banner_app
// @Summary      Get all banner_app
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @Param 		 type   query     string  false "Type of banner (e.g., 'anime', 'movie', kdrama)"
// @Router       /app/banner [get]
func (h *BannerAppController) GetAllBannerApp(c *fiber.Ctx) error {
	query := &request.QueryBannerApp{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 10),
		Type:  c.Query("type"),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.BannerApp]{
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

// @Tags         banner_app
// @Summary      Get a banner_app by ID
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "BannerApp ID (uint)"
// @Router       /app/banner/{id} [get]
func (h *BannerAppController) GetBannerAppByID(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.BannerApp]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         banner_app
// @Summary      Create a new banner_app
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  request.CreateBannerApp  true  "Request body"
// @Router       /app/banner [post]
func (h *BannerAppController) CreateBannerApp(c *fiber.Ctx) error {
	req := new(request.CreateBannerApp)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.BannerApp]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         banner_app
// @Summary      Update a banner_app
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "BannerApp ID (uint)"
// @Param        request  body  request.UpdateBannerApp  true  "Request body"
// @Router       /app/banner/{id} [put]
func (h *BannerAppController) UpdateBannerApp(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	req := new(request.UpdateBannerApp)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.BannerApp]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         banner_app
// @Summary      Delete a banner_app
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "BannerApp ID (uint)"
// @Router       /app/banner/{id} [delete]
func (h *BannerAppController) DeleteBannerApp(c *fiber.Ctx) error {
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

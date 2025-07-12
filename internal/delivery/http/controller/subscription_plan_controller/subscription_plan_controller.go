package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/subscription_plan/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/subscription_plan_service"
)

type SubscriptionPlanController struct {
	Service service.SubscriptionPlanServiceInterface
}

func NewSubscriptionPlanController(service service.SubscriptionPlanServiceInterface) *SubscriptionPlanController {
	return &SubscriptionPlanController{Service: service}
}

// @Tags         subscription_plan
// @Summary      Get all subscription_plan
// @Produce      json
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @Router        /subscription/plans [get]
func (h *SubscriptionPlanController) GetAllSubscriptionPlan(c *fiber.Ctx) error {
	query := &request.QuerySubscriptionPlan{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 10),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.SubscriptionPlan]{
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

// @Tags         subscription_plan
// @Summary      Get a subscription_plan by ID
// @Produce      json
// @Param        id  path  int  true  "SubscriptionPlan ID (uint)"
// @Router        /subscription/plans/{id} [get]
func (h *SubscriptionPlanController) GetSubscriptionPlanByID(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.SubscriptionPlan]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         subscription_plan
// @Summary      Create a new subscription_plan
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  request.CreateSubscriptionPlan  true  "Request body"
// @Router        /subscription/plans [post]
func (h *SubscriptionPlanController) CreateSubscriptionPlan(c *fiber.Ctx) error {
	req := new(request.CreateSubscriptionPlan)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.SubscriptionPlan]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         subscription_plan
// @Summary      Update a subscription_plan
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "SubscriptionPlan ID (uint)"
// @Param        request  body  request.UpdateSubscriptionPlan  true  "Request body"
// @Router        /subscription/plans/{id} [put]
func (h *SubscriptionPlanController) UpdateSubscriptionPlan(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	req := new(request.UpdateSubscriptionPlan)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.SubscriptionPlan]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         subscription_plan
// @Summary      Delete a subscription_plan
// @Security 	 BearerAuth
// @Produce      json
// @Param        id  path  int  true  "SubscriptionPlan ID (uint)"
// @Router        /subscription/plans/{id} [delete]
func (h *SubscriptionPlanController) DeleteSubscriptionPlan(c *fiber.Ctx) error {
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

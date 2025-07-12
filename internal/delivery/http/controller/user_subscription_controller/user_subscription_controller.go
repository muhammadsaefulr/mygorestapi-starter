package controller

import (
	"math"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_subscription/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_subscription_service"
)

type UserSubscriptionController struct {
	Service service.UserSubscriptionServiceInterface
}

func NewUserSubscriptionController(service service.UserSubscriptionServiceInterface) *UserSubscriptionController {
	return &UserSubscriptionController{Service: service}
}

// @Tags         user_subscription
// @Summary      Get all user_subscription
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @Router       /user/subscriptions [get]
func (h *UserSubscriptionController) GetAllUserSubscription(c *fiber.Ctx) error {
	query := &request.QueryUserSubscription{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 10),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.UserSubscription]{
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

// @Tags         user_subscription
// @Summary      Get a user_subscription by ID
// @Security     BearerAuth
// @Produce      json
// @Param        user_id  path  string  true  "User ID (string)"
// @Router       /user/subscriptions/{user_id} [get]
func (h *UserSubscriptionController) GetUserSubscriptionByID(c *fiber.Ctx) error {
	idStr := c.Params("user_id")

	result, err := h.Service.GetByUserID(c, idStr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.UserSubscription]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         user_subscription
// @Summary      Create a new user_subscription
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  request.CreateUserSubscription  true  "Request body"
// @Router       /user/subscriptions [post]
func (h *UserSubscriptionController) CreateUserSubscription(c *fiber.Ctx) error {
	req := new(request.CreateUserSubscription)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.UserSubscription]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         user_subscription
// @Summary      Update a user_subscription
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user_id       path  string  true  "User ID (string)"
// @Param        request  body  request.UpdateUserSubscription  true  "Request body"
// @Router       /user/subscriptions/{user_id} [put]
func (h *UserSubscriptionController) UpdateUserSubscription(c *fiber.Ctx) error {
	idStr := c.Params("user_id")

	req := new(request.UpdateUserSubscription)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.UpdateByUserId(c, idStr, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.UserSubscription]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         user_subscription
// @Summary      Delete a user_subscription
// @Security     BearerAuth
// @Produce      json
// @Param        user_id  path  string  true  "User ID (string)"
// @Router       /user/subscriptions/{user_id} [delete]
func (h *UserSubscriptionController) DeleteUserSubscription(c *fiber.Ctx) error {
	idStr := c.Params("user_id")

	if err := h.Service.DeleteByUserId(c, idStr); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data deleted successfully",
	})
}

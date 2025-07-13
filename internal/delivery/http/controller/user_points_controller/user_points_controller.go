package controller

import (
	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_points/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_points_service"
)

type UserPointsController struct {
	Service service.UserPointsService
}

func NewUserPointsController(service service.UserPointsService) *UserPointsController {
	return &UserPointsController{Service: service}
}

// @Tags         user_points
// @Summary      Get a user_points by ID
// @Produce      json
// @Param        id  path  string  true  "UserPoints User ID (string)"
// @Router       /users/points/{id} [get]
func (h *UserPointsController) GetUserPointsByID(c *fiber.Ctx) error {
	idStr := c.Params("id")

	result, err := h.Service.GetByUserID(c, idStr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.UserPoints]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         user_points
// @Summary      Create a new user_points
// @Security BearerAuth
// @Description Create a new user_points, type update can be 'add' or 'subtract'
// @Accept       json
// @Produce      json
// @Param        request  body  request.UserPoints  true  "Request body"
// @Router      /users/points [post]
func (h *UserPointsController) PostUserPoints(c *fiber.Ctx) error {
	req := new(request.UserPoints)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Post(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.UserPoints]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         user_points
// @Summary      Delete a user_points
// @Produce      json
// @Security 	 BearerAuth
// @Param        id  path  int  true  "UserPoints ID (uint)"
// @Router      /users/points/{id} [delete]
func (h *UserPointsController) DeleteUserPoints(c *fiber.Ctx) error {
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

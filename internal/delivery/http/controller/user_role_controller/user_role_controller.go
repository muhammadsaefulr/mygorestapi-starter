package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_role/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_role_service"
)

type UserRoleController struct {
	Service service.UserRoleService
}

func NewUserRoleController(service service.UserRoleService) *UserRoleController {
	return &UserRoleController{Service: service}
}

// @Tags         user_role
// @Summary      Get all user_role
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @Router       /user/role [get]
func (h *UserRoleController) GetAllUserRole(c *fiber.Ctx) error {
	query := &request.QueryUserRole{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 10),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.UserRole]{
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

// @Tags         user_role
// @Summary      Get a user_role by ID
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "UserRole ID (uint)"
// @Router       /user/role/{id} [get]
func (h *UserRoleController) GetUserRoleByID(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.UserRole]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         user_role
// @Summary      Create a new user_role
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  request.CreateUserRole  true  "Request body"
// @Router       /user/role [post]
func (h *UserRoleController) CreateUserRole(c *fiber.Ctx) error {
	req := new(request.CreateUserRole)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.UserRole]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         user_role
// @Summary      Update a user_role
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "UserRole ID (uint)"
// @Param        request  body  request.UpdateUserRole  true  "Request body"
// @Router       /user/role/{id} [put]
func (h *UserRoleController) UpdateUserRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	req := new(request.UpdateUserRole)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.UserRole]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         user_role
// @Summary      Delete a user_role
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "UserRole ID (uint)"
// @Router       /user/role/{id} [delete]
func (h *UserRoleController) DeleteUserRole(c *fiber.Ctx) error {
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

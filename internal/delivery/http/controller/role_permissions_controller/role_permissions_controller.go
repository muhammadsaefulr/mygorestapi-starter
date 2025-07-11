package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/role_permissions/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/role_permissions_service"
)

type RolePermissionsController struct {
	Service service.RolePermissionsService
}

func NewRolePermissionsController(service service.RolePermissionsService) *RolePermissionsController {
	return &RolePermissionsController{Service: service}
}

// @Tags         role_permissions
// @Summary      Get all role_permissions
// @Produce      json
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @Router       /user/roles/permissions [get]
func (h *RolePermissionsController) GetAllRolePermissions(c *fiber.Ctx) error {
	query := &request.QueryRolePermissions{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 10),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.RolePermissions]{
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

// @Tags         role_permissions
// @Summary      Get a role_permissions by ID
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "RolePermissions ID (uint)"
// @Router       /user/roles/permissions/{id} [get]
func (h *RolePermissionsController) GetRolePermissionsByID(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.RolePermissions]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         role_permissions
// @Summary      Create a new role_permissions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  request.CreateRolePermissions  true  "Request body"
// @Router       /user/roles/permissions [post]
func (h *RolePermissionsController) CreateRolePermissions(c *fiber.Ctx) error {
	req := new(request.CreateRolePermissions)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.RolePermissions]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         role_permissions
// @Summary      Update a role_permissions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "RolePermissions ID (uint)"
// @Param        request  body  request.UpdateRolePermissions  true  "Request body"
// @Router       /user/roles/permissions/{id} [put]
func (h *RolePermissionsController) UpdateRolePermissions(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	req := new(request.UpdateRolePermissions)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.RolePermissions]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         role_permissions
// @Summary      Delete a role_permissions
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "RolePermissions ID (uint)"
// @Router       /user/roles/permissions/{id} [delete]
func (h *RolePermissionsController) DeleteRolePermissions(c *fiber.Ctx) error {
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

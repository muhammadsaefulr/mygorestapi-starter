package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_badge/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_badge_service"
)

type UserBadgeController struct {
	Service service.UserBadgeService
}

func NewUserBadgeController(service service.UserBadgeService) *UserBadgeController {
	return &UserBadgeController{Service: service}
}

// @Tags         user_badge
// @Summary      Get all user_badge
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @Router       /badge [get]
func (h *UserBadgeController) GetAllUserBadge(c *fiber.Ctx) error {
	query := &request.QueryUserBadge{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 10),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.UserBadge]{
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

// @Tags         user_badge
// @Summary      Get a user_badge by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "UserBadge ID (uint)"
// @Router       /badge/{id} [get]
func (h *UserBadgeController) GetUserBadgeByID(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.UserBadge]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         user_badge
// @Summary      Create a new user_badge
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  request.CreateUserBadge  true  "Request body"
// @Router       /badge [post]
func (h *UserBadgeController) CreateUserBadge(c *fiber.Ctx) error {
	req := new(request.CreateUserBadge)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.UserBadge]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         user_badge
// @Summary      Update a user_badge
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "UserBadge ID (uint)"
// @Param        request  body  request.UpdateUserBadge  true  "Request body"
// @Router       /badge/{id} [put]
func (h *UserBadgeController) UpdateUserBadge(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	req := new(request.UpdateUserBadge)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.UserBadge]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         user_badge
// @Summary      Delete a user_badge
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  int  true  "UserBadge ID (uint)"
// @Router       /badge/{id} [delete]
func (h *UserBadgeController) DeleteUserBadge(c *fiber.Ctx) error {
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

// @Tags         user_badge
// @Summary      Get user badge info by user id
// @Security     BearerAuth
// @Produce      json
// @Param        user_id  path  string  true  "User ID"
// @Router       /user/badge/{user_id} [get]
func (h *UserBadgeController) GetUserBadgeInfoByUserID(c *fiber.Ctx) error {
	user_id := c.Params("user_id")

	data, err := h.Service.GetUserBadgeInfoByUserID(c, user_id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[model.UserBadgeInfo]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Results: data,
	})
}

// @Tags         user_badge
// @Summary      Create user badge info
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  request.CreateUserBadgeInfo  true  "Request body"
// @Router       /user/badge [post]
func (h *UserBadgeController) CreateUserBadgeInfo(c *fiber.Ctx) error {
	req := new(request.CreateUserBadgeInfo)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	err := h.Service.CreateUserBadgeInfo(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.Common{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
	})
}

// @Tags         user_badge
// @Summary      Update user badge info
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  request.UpdateUserBadgeInfo  true  "Request body"
// @Router       /user/badge/{user_id} [put]
func (h *UserBadgeController) UpdateUserBadgeInfo(c *fiber.Ctx) error {
	user_id := c.Params("user_id")

	req := new(request.UpdateUserBadgeInfo)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	req.UserID = user_id

	err := h.Service.UpdateUserBadgeInfo(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
	})
}

// @Tags         user_badge
// @Summary      Delete user badge info
// @Security     BearerAuth
// @Produce      json
// @Param        user_id  path  string  true  "User ID"
// @Router       /user/badge/{user_id} [delete]
func (h *UserBadgeController) DeleteUserBadgeInfo(c *fiber.Ctx) error {
	user_id := c.Params("user_id")

	err := h.Service.DeleteUserBadgeInfo(c, user_id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data deleted successfully",
	})
}

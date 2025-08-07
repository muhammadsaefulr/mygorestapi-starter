package controller

import (
	"math"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"

	responses "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/response"
	system_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/system_service"
	user_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController struct {
	UserService  user_service.UserService
	TokenService system_service.TokenService
}

func NewUserController(userService user_service.UserService, tokenService system_service.TokenService) *UserController {
	return &UserController{
		UserService:  userService,
		TokenService: tokenService,
	}
}

// @Tags         Users
// @Summary      Get total users
// @Produce      json
// @Router       /users/count [get]
func (u *UserController) GetCountAllUser(c *fiber.Ctx) error {
	total, err := u.UserService.GetCountAllUser(c)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[int64]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Get total users successfully",
		Data:    total,
	})
}

// @Tags         Users
// @Summary      Get all users
// @Description  Only admins can retrieve all users.
// @Security BearerAuth
// @Produce      json
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of users"    default(10)
// @Param        search   query     string  false  "Search by name or email or role"
// @Param        role     query     string  false  "Search by role" (default "1")
// @Router       /users [get]
// @Success      200  {object}  example.GetAllUserResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
func (u *UserController) GetUsers(c *fiber.Ctx) error {
	query := &request.QueryUser{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Role:   c.Query("role", ""),
		Search: c.Query("search", ""),
	}

	users, totalResults, err := u.UserService.GetAllUser(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithPaginate[model.User]{
			Code:         fiber.StatusOK,
			Status:       "success",
			Message:      "Get all users successfully",
			Results:      users,
			Page:         query.Page,
			Limit:        query.Limit,
			TotalPages:   int64(math.Ceil(float64(totalResults) / float64(query.Limit))),
			TotalResults: totalResults,
		})
}

// @Tags         Users
// @Summary      Get user session
// @Description  Get user session
// @Security BearerAuth
// @Produce      json
// @Router       /users/session [get]
// @Success      200  {object}  responses.GetUsersResponse
func (u *UserController) GetUserSession(c *fiber.Ctx) error {
	user, err := u.UserService.GetUserSession(c)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[responses.GetUsersResponse]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Get user session successfully",
		Data:    *convert_types.UserModelToUserResponse(user),
	})
}

// @Tags         Users
// @Summary      Get a user
// @Description  Logged in users can fetch only their own user information. Only admins can fetch other users.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "User id"
// @Router       /users/{id} [get]
// @Success      200  {object}  example.GetUserResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (u *UserController) GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("userId")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := u.UserService.GetUserByID(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithUser{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Get user successfully",
			User:    *user,
		})
}

// @Tags         Users
// @Summary      Create a user
// @Description  Only admins can create other users.
// @Security BearerAuth
// @Produce      json
// @Param        request  body  request.CreateUser  true  "Request body"
// @Router       /users [post]
// @Success      201  {object}  example.CreateUserResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      409  {object}  example.DuplicateEmail  "Email already taken"
func (u *UserController) CreateUser(c *fiber.Ctx) error {
	req := new(request.CreateUser)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := u.UserService.CreateUser(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithUser{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Create user successfully",
			User:    *user,
		})
}

// @Tags         Users
// @Summary      Update a user
// @Description  Logged in users can only update their own information. Only admins can update other users.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "User id"
// @Param        request  body  request.UpdateUser  true  "Request body"
// @Router       /users/{id} [put]
// @Success      200  {object}  example.UpdateUserResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
// @Failure      409  {object}  example.DuplicateEmail  "Email already taken"
func (u *UserController) UpdateUser(c *fiber.Ctx) error {
	req := new(request.UpdateUser)
	userID := c.Params("userId")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	res, err := u.UserService.UpdateUser(c, userID, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithUser{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Update user successfully",
			User:    *res,
		})
}

// @Tags         Users
// @Summary      Delete a user
// @Description  Logged in users can delete only themselves. Only admins can delete other users.
// @Security BearerAuth
// @Produce      json
// @Param        id  path  string  true  "User id"
// @Router       /users/{id} [delete]
// @Success      200  {object}  example.DeleteUserResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
// @Failure      403  {object}  example.Forbidden  "Forbidden"
// @Failure      404  {object}  example.NotFound  "Not found"
func (u *UserController) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("userId")

	if _, err := uuid.Parse(userID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := u.TokenService.DeleteAllToken(c, userID); err != nil {
		return err
	}

	if err := u.UserService.DeleteUser(c, userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Delete user successfully",
		})
}

package controller

import (
	"math"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/request"
	responses "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/movie_details_service"
)

type MovieDetailsController struct {
	Service service.MovieDetailsServiceInterface
}

func NewMovieDetailsController(service service.MovieDetailsServiceInterface) *MovieDetailsController {
	return &MovieDetailsController{Service: service}
}

// @Tags         movie
// @Summary      Get all movie_details
// @Produce      json
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term like title or studio"
// @Param        type   query     string  false  "Type of movie" default(anime)
// @success      200    {object}  response.SuccessWithPaginate[model.MovieDetails]  "Successfully retrieved data"
// @Router       /movie/details [get]
func (h *MovieDetailsController) GetAllMovieDetails(c *fiber.Ctx) error {
	query := &request.QueryMovieDetails{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Type:   c.Query("type"),
		Search: c.Query("search"),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.MovieDetails]{
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

// @Tags         movie
// @Summary      Get a movie_details by ID
// @Produce      json
// @Param        id  path  string  true  "MovieDetails ID (uint)"
// @success      200    {object}  response.SuccessWithDetail[model.MovieDetails]  "Data retrieved successfully"
// @Router       /movie/details/{id} [get]
func (h *MovieDetailsController) GetMovieDetailsByID(c *fiber.Ctx) error {
	idStr := c.Params("id")

	result, err := h.Service.GetByIDPreEps(c, idStr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[responses.MovieDetailsResponse]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         movie
// @Summary      Create a new movie_details
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  request.CreateMovieDetails  true  "Request body"
// @success      201    {object}  response.SuccessWithDetail[model.MovieDetails]  "Data created successfully"
// @Router       /movie/details [post]
func (h *MovieDetailsController) CreateMovieDetails(c *fiber.Ctx) error {
	req := new(request.CreateMovieDetails)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.MovieDetails]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         movie
// @Summary      Update a movie_details
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path  string  true  "MovieDetails ID (uint)"
// @Param        request  body  request.UpdateMovieDetails  true  "Request body"
// @success      200    {object}  response.SuccessWithDetail[model.MovieDetails]  "Data updated successfully"
// @Router       /movie/details/{id} [put]
func (h *MovieDetailsController) UpdateMovieDetails(c *fiber.Ctx) error {
	idStr := c.Params("id")

	req := new(request.UpdateMovieDetails)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, idStr, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.MovieDetails]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         movie
// @Summary      Delete a movie_details
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "MovieDetails ID (uint)"
// @Router       /movie/details/{id} [delete]
func (h *MovieDetailsController) DeleteMovieDetails(c *fiber.Ctx) error {
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

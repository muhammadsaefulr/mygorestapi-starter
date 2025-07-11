package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_movie/request"
	responseReqMov "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_movie/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/request_movie_service"
)

type RequestMovieController struct {
	Service service.RequestMovieService
}

func NewRequestMovieController(service service.RequestMovieService) *RequestMovieController {
	return &RequestMovieController{Service: service}
}

// @Tags         RequestMovie
// @Summary      Get all request_movie
// @Produce      json
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @Security     BearerAuth
// @Success      200  {object}    response.SuccessWithPaginate[responseReqMov.RequestMovieResponse]
// @Router       /request-movie [get]
func (h *RequestMovieController) GetAllRequestMovie(c *fiber.Ctx) error {
	query := &request.QueryRequestMovie{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 10),
	}

	data, total, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[responseReqMov.RequestMovieResponse]{
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

// @Tags         RequestMovie
// @Summary      Get a request_movie by ID
// @Produce      json
// @Param        id  path  int  true  "RequestMovie ID (uint)"
// @Security     BearerAuth
// @Success      200  {object}    response.SuccessWithDetail[responseReqMov.RequestMovieResponse]
// @Router       /request-movie/{id} [get]
func (h *RequestMovieController) GetRequestMovieByID(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[responseReqMov.RequestMovieResponse]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         RequestMovie
// @Summary      Create a new request_movie
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  request.CreateRequestMovie  true  "Request body"
// @Router       /request-movie [post]
func (h *RequestMovieController) CreateRequestMovie(c *fiber.Ctx) error {
	req := new(request.CreateRequestMovie)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.RequestMovie]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         RequestMovie
// @Summary      Update a request_movie
// @Accept       json
// @Produce      json
// @Param        id       path  int  true  "RequestMovie ID (uint)"
// @Security     BearerAuth
// @Param        request  body  request.UpdateRequestMovie  true  "Request body"
// @Router       /request-movie/{id} [put]
func (h *RequestMovieController) UpdateRequestMovie(c *fiber.Ctx) error {
	idStr := c.Params("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format, must be a positive integer")
	}
	id := uint(idVal)

	req := new(request.UpdateRequestMovie)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[responseReqMov.RequestMovieResponse]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         RequestMovie
// @Summary      Delete a request_movie
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "RequestMovie ID (uint)"
// @Router       /request-movie/{id} [delete]
func (h *RequestMovieController) DeleteRequestMovie(c *fiber.Ctx) error {
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

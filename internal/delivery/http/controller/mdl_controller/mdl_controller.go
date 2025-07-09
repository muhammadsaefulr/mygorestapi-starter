package controller

import (
	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/mdl/request"
	responses "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/mdl_service"
)

type MdlController struct {
	Service service.MdlServiceInterface
}

func NewMdlController(service service.MdlServiceInterface) *MdlController {
	return &MdlController{Service: service}
}

// @Tags         Mdl
// @Summary      Get all mdl
// @Produce      json
// @Param        page   query     int     false  "Page number"  default(1)
// @Param        limit  query     int     false  "Items per page"  default(10)
// @Param        search query     string  false  "Search term"
// @Param        category  query     string  false  "Discovery category"  Enums(popular, trending, ongoing, rekom)  default(popular)
// @Router       /mdl [get]
func (h *MdlController) GetAllMdl(c *fiber.Ctx) error {
	query := &request.QueryMdl{
		Page:     c.QueryInt("page", 1),
		Limit:    c.QueryInt("limit", 10),
		Category: c.Query("category", "trending"),
	}

	data, total, totalPages, err := h.Service.GetAll(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[responses.MovieDetailOnlyResponse]{
		Code:         fiber.StatusOK,
		Status:       "success",
		Message:      "Successfully retrieved data",
		Results:      data,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: total,
	})
}

// @Tags         Mdl
// @Summary      Get a movie_details by ID
// @Produce      json
// @Param        id  path  string  true  "MovieDetails ID (uint)"
// @success      200    {object}  response.SuccessWithDetail[responses.MovieDetailOnlyResponse]  "Data retrieved successfully"
// @Router       /mdl/{id} [get]
func (h *MdlController) GetMdlByID(c *fiber.Ctx) error {
	idStr := c.Params("id")

	result, err := h.Service.GetDetailByID(c, idStr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[responses.MovieDetailOnlyResponse]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

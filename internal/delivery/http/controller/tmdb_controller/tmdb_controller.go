package controller

import (
	"math"

	"github.com/gofiber/fiber/v2"

	responses "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/tmdb/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/tmdb_service"
)

type TmdbController struct {
	Service service.TmdbServiceInterface
}

func NewTmdbController(service service.TmdbServiceInterface) *TmdbController {
	return &TmdbController{Service: service}
}

// @Tags         Tmdb
// @Summary      Get list of TMDb media (movie)
// @Description  Retrieve TMDb-based data (movie) by category (popular, trending, etc).
// @Produce      json
// @Param        page     query     int     false  "Page number"       default(1)
// @Param        limit    query     int     false  "Items per page, max limit 20"    default(10)
// @Param        type     query     string  false  "Media type"        Enums(tv, movie) default(movie)
// @Param        category query     string  true   "Media category"    Enums(popular, trending, ongoing, rekom)
// @Param        search   query     string  false  "Search keyword (for rekom category need this param)"
// @Success      200      {object}  response.SuccessWithPaginate[responses.MovieDetailOnlyResponse]
// @Router       /tmdb [get]
func (h *TmdbController) GetAllTmdb(c *fiber.Ctx) error {
	query := &request.QueryTmdb{
		Page:     c.QueryInt("page", 1),
		Limit:    c.QueryInt("limit", 10),
		Type:     c.Query("type", "movie"),
		Category: c.Query("category", ""),
		Search:   c.Query("search", ""),
	}

	data, total, err := h.Service.GetAll(c, query)
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
		TotalPages:   int64(math.Ceil(float64(total) / float64(query.Limit))),
		TotalResults: total,
	})
}

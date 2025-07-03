package controller

import (
	"math"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/anilist/request"
	responses "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/anilist_service"
)

type AnilistController struct {
	Service service.AnilistServiceInterface
}

func NewAnilistController(service service.AnilistServiceInterface) *AnilistController {
	return &AnilistController{Service: service}
}

// @Tags         Anilist
// @Summary      Get anime discovery list (popular, trending, ongoing, rekomendasi)
// @Description  Retrieve anime from AniList API based on category. When using category=rekom, 'search' is required to perform a title-based recommendation.
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)        minimum(1)
// @Param        limit     query     int     false  "Items per page"  default(10)    minimum(1)  maximum(50)
// @Param        category  query     string  false  "Discovery category"  Enums(popular, trending, ongoing, rekom)  default(popular)
// @Param        search    query     string  false  "Search term. Required if category is 'rekom'"  default(one piece)
// @Success      200       {object}  response.SuccessWithPaginate[model.MovieDetails]
// @Router       /anilists [get]
func (h *AnilistController) GetAllAnilist(c *fiber.Ctx) error {
	query := new(request.QueryAnilist)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid query params",
		})
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

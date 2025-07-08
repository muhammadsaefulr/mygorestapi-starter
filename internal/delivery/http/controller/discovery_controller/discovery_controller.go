package controller

import (
	"log"
	"math"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/discovery/request"
	responses "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/discovery_service"
)

type DiscoveryController struct {
	Service service.DiscoveryServiceInterface
}

func NewDiscoveryController(service service.DiscoveryServiceInterface) *DiscoveryController {
	return &DiscoveryController{Service: service}
}

// @Tags         discovery
// @Summary      Get popular discover
// @Description  Get popular discover
// @Accept       json
// @Produce      json
// @Param        page query     int     false "Page"
// @Param        limit query    int     false "Limit"
// @Param        type query     string  false  "Movie Type"  Enums(anime, kdrama, tv, movie)  default(anime)
// @Param        category query string  false  "Discovery Category"  Enums(popular, trending, ongoing)  default(popular)
// @Router       /discovery/popular [get]
func (c *DiscoveryController) GetDiscover(ctx *fiber.Ctx) error {
	params := &request.QueryDiscovery{}

	if err := ctx.QueryParser(params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query params")
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	data, total, err := c.Service.GetDiscover(ctx, params)
	if err != nil {
		log.Printf("GetPopularDiscover error: %v", err) // ✅ Tambah ini
		return err
	}

	return ctx.JSON(response.SuccessWithPaginate[responses.MovieDetailOnlyResponse]{
		Code:         fiber.StatusOK,
		Status:       "success",
		Message:      "Successfully retrieved data",
		Results:      data,
		Page:         params.Page,
		Limit:        params.Limit,
		TotalPages:   int64(math.Ceil(float64(total) / float64(params.Limit))),
		TotalResults: total,
	})
}

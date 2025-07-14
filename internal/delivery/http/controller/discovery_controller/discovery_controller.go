package controller

import (
	"log"

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
// @Param        type query     string  false "Movie Type"  Enums(anime, kdrama, tv, movie)  default(anime)
// @Param        genre query    string  false "Genre"
// @Param        category query string  false  "Discovery Category"  Enums(popular, trending, ongoing, genre, search)  default(popular)
// @Param        search query string false "Search term (Only used when category=search)"
// @Router       /discovery [get]
func (c *DiscoveryController) GetDiscover(ctx *fiber.Ctx) error {
	params := &request.QueryDiscovery{
		Page:     ctx.QueryInt("page", 1),
		Limit:    ctx.QueryInt("limit", 10),
		Category: ctx.Query("category"),
		Search:   ctx.Query("search"),
		Genre:    ctx.Query("genre"),
		Type:     ctx.Query("type"),
	}

	if err := ctx.QueryParser(params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query params")
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	data, page, total, err := c.Service.GetDiscover(ctx, params)
	if err != nil {
		log.Printf("GetPopularDiscover error: %v", err)
		return err
	}

	return ctx.JSON(response.SuccessWithPaginate[responses.MovieDetailOnlyResponse]{
		Code:         fiber.StatusOK,
		Status:       "success",
		Message:      "Successfully retrieved data",
		Results:      data,
		Page:         params.Page,
		Limit:        params.Limit,
		TotalPages:   page,
		TotalResults: total,
	})
}

// @Tags         discovery
// @Summary      Get detail by title
// @Description  Get detail by title
// @Accept       json
// @Produce      json
// @Param        title path     string  true "Title"
// @Param        mediaType path string  true "Media Type"  Enums(anime, kdrama, tv, movie)  default(anime)
// @Router       /discovery/detail/{mediaType}/{title} [get]
func (c *DiscoveryController) GetDiscoverDetailByTitle(ctx *fiber.Ctx) error {
	mediaType := ctx.Params("mediaType")
	title := ctx.Params("title")

	result, err := c.Service.GetDiscoverDetailByTitle(ctx, mediaType, title)
	if err != nil {
		return err
	}

	return ctx.JSON(response.SuccessWithDetail[responses.MovieDetailOnlyResponse]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

package controller

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/watchlist_service"
)

type WatchlistController struct {
	WatchlistService service.WatchlistService
}

func NewWatchlistController(service service.WatchlistService) *WatchlistController {
	return &WatchlistController{
		WatchlistService: service,
	}
}

// @Tags         Watchlist
// @Summary      Get all watchlists
// @Description  Only admins can retrieve all watchlists.
// @Security BearerAuth
// @Produce      json
// @Router       /watchlists/ [get]
func (h *WatchlistController) GetAllWatchlist(c *fiber.Ctx) error {
	query := &request.QueryWatchlist{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	watchlists, totalResults, err := h.WatchlistService.GetAllWatchlist(c, query)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithPaginate[model.Watchlist]{
			Code:         fiber.StatusOK,
			Status:       "success",
			Message:      "Get all watchlists successfully",
			Results:      watchlists,
			Page:         query.Page,
			Limit:        query.Limit,
			TotalPages:   int64(math.Ceil(float64(totalResults) / float64(query.Limit))),
			TotalResults: totalResults,
		})
}

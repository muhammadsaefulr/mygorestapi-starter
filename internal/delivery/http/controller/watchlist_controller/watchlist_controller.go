package controller

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/request"
	responseWatchlist "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/response"
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
		JSON(response.SuccessWithPaginate[responseWatchlist.WatchlistResponse]{
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

// @Tags         Watchlist
// @Summary      Create watchlist
// @Description  User Create watchlist
// @Security BearerAuth
// @Produce      json
// @Param        request  body  request.CreateWatchlist  true  "Request body"
// @Router       /watchlists [post]
func (h *WatchlistController) CreateWatchlist(c *fiber.Ctx) error {
	req := new(request.CreateWatchlist)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	watchlist, err := h.WatchlistService.CreateWatchlist(c, req)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create watchlist")
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.Watchlist]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Create watchlist successfully",
		Data: model.Watchlist{
			UserId:    watchlist.UserId,
			MovieId:   watchlist.MovieId,
			ID:        watchlist.ID,
			CreatedAt: watchlist.CreatedAt,
			UpdatedAt: watchlist.UpdatedAt,
		},
	})

}

// @Tags         Watchlist
// @Summary      Update watchlist
// @Description  User Update watchlist
// @Security BearerAuth
// @Param        id  path  string  true  "Watchlist id"
// @Param        request  body  request.UpdateWatchlist  true  "Request body"
// @Produce      json
// @Router       /watchlists/{id} [put]
func (h *WatchlistController) UpdateWatchlist(c *fiber.Ctx) error {
	req := new(request.UpdateWatchlist)
	id := c.Params("id")

	uId, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid id")
	}

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	watchlist, err := h.WatchlistService.UpdateWatchlist(c, uint(uId), req)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update watchlist")
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.Watchlist]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Update watchlist successfully",
		Data: model.Watchlist{
			UserId:    watchlist.UserId,
			MovieId:   watchlist.MovieId,
			ID:        watchlist.ID,
			CreatedAt: watchlist.CreatedAt,
			UpdatedAt: watchlist.UpdatedAt,
		},
	})
}

// @Tags         Watchlist
// @Summary      Delete watchlist
// @Description  User Delete watchlist
// @Security BearerAuth
// @Produce      json
// @Router       /watchlists/{id} [delete]
func (h *WatchlistController) DeleteWatchlist(c *fiber.Ctx) error {
	id := c.Params("id")

	uId, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid id")
	}

	if err := h.WatchlistService.DeleteWatchlist(c, uint(uId)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete watchlist")
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Delete watchlist successfully",
	})
}

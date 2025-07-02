package controller

import (
	// "math"

	"github.com/gofiber/fiber/v2"

	request "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/request"
	responses "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/movie_episode_service"
)

type MovieEpisodeController struct {
	Service service.MovieEpisodeServiceInterface
}

func NewMovieEpisodeController(service service.MovieEpisodeServiceInterface) *MovieEpisodeController {
	return &MovieEpisodeController{Service: service}
}

// // @Tags         movie
// // @Summary      Get all movie episodes
// // @Produce      json
// // @Param        page   query     int     false  "Page number"  default(1)
// // @Param        limit  query     int     false  "Items per page"  default(10)
// // @Param        search query     string  false  "Search term"
// // @Router       /movie/episodes [get]
// func (h *MovieEpisodeController) GetAllMovieEpisode(c *fiber.Ctx) error {
// 	query := &request.QueryMovieEpisode{
// 		Page:  c.QueryInt("page", 1),
// 		Limit: c.QueryInt("limit", 10),
// 	}

// 	data, total, err := h.Service.GetAll(c, query)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.MovieEpisode]{
// 		Code:         fiber.StatusOK,
// 		Status:       "success",
// 		Message:      "Successfully retrieved data",
// 		Results:      data,
// 		Page:         query.Page,
// 		Limit:        query.Limit,
// 		TotalPages:   int64(math.Ceil(float64(total) / float64(query.Limit))),
// 		TotalResults: total,
// 	})
// }

// @Tags         movie
// @Summary      Get a movie episodes by ID
// @Produce      json
// @Param        movie_eps_id  path  string  true  "MovieEps ID (string)"
// @Router       /movie/episodes/{movie_eps_id} [get]
func (h *MovieEpisodeController) GetMovieEpisodeByID(c *fiber.Ctx) error {
	idStr := c.Params("id")

	result, err := h.Service.GetByID(c, idStr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.MovieEpisode]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         movie
// @Summary      Get all movie episodes by movie ID
// @Produce      json
// @Param        movie_id  path  string  true  "Movie ID"
// @Param        movie_eps_id  path  string  true  "Movie Eps ID"
// @Router       /movie/episodes/{movie_id}/{movie_eps_id} [get]
func (h *MovieEpisodeController) GetMovieEpisodeByMovieID(c *fiber.Ctx) error {
	idStr := c.Params("movie_id")
	movEpsId := c.Params("movie_eps_id")

	result, err := h.Service.GetByMovieID(c, movEpsId, idStr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[responses.MovieEpisodeResponses]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    *result,
	})
}

// @Tags         movie
// @Summary      Create a new movie episodes
// @Accept       application/json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  request.CreateMovieEpisodes  true  "Request body"
// @Success      201 {object} response.SuccessWithDetail[model.MovieEpisode]
// @Router       /movie/episodes [post]
func (h *MovieEpisodeController) CreateMovieEpisodes(c *fiber.Ctx) error {
	req := new(request.CreateMovieEpisodes)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Create(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.MovieEpisode]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Data created successfully",
		Data:    *result,
	})
}

// @Tags         movie
// @Summary      Upload new movie episode (with file)
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        movie_eps_id   formData string true  "Episode ID"
// @Param        movie_id       formData string true  "Movie ID"
// @Param        title          formData string true  "Title"
// @Param        resolution     formData string true  "Resolution"
// @Param        file_video     formData file   true  "Video File"
// @Success      201 {object} response.SuccessWithDetail[model.MovieEpisode]
// @Router       /movie/episodes/upload [post]
func (h *MovieEpisodeController) CreateUpload(c *fiber.Ctx) error {
	req := new(request.CreateMovieEpisodesUpload)

	fileHeader, err := c.FormFile("file_video")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Video file is required")
	}
	req.ContentUploads = fileHeader

	result, err := h.Service.CreateUpload(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithDetail[model.MovieEpisode]{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Episode uploaded successfully",
		Data:    *result,
	})
}

// @Tags         movie
// @Summary      Update a movie episodes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        movie_eps_id  path  string  true  "MovieEps ID (string)"
// @Param        request  body  request.UpdateMovieEpisodes  true  "Request body"
// @Router       /movie/episodes/{id} [put]
func (h *MovieEpisodeController) UpdateMovieEpisodes(c *fiber.Ctx) error {
	idStr := c.Params("id")

	req := new(request.UpdateMovieEpisodes)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.Service.Update(c, idStr, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.MovieEpisode]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Data updated successfully",
		Data:    *result,
	})
}

// @Tags         movie
// @Summary      Delete a movie episodes
// @Produce      json
// @Security     BearerAuth
// @Param        movie_eps_id  path  string  true  "MovieEps ID (string)"
// @Router       /movie/episodes/{movie_eps_id} [delete]
func (h *MovieEpisodeController) DeleteMovieEpisode(c *fiber.Ctx) error {
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

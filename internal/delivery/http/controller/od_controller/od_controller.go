package controller

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/util/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"

	od_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/otakudesu_scrape"

	"github.com/gofiber/fiber/v2"
)

type OdAnimeController struct {
	AnimeService od_service.AnimeService
}

func NewAnimeController(animeService od_service.AnimeService) *OdAnimeController {
	return &OdAnimeController{
		AnimeService: animeService,
	}
}

// @Tags         Otakudesu
// @Summary      Get homepage anime data
// @Description  Scrape and get list of anime from Otakudesu homepage.
// @Produce      json
// @Router       /otakudesu/ [get]
// @Success      200  {object}  example.GetOdAnimeHomeResponse
func (a *OdAnimeController) GetHomePageAnime(c *fiber.Ctx) error {
	animes, err := a.AnimeService.GetHomePage()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[model.AnimeData]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Successfully Retrieved Anime!",
		Results: animes,
	})
}

// @Tags         Otakudesu
// @Summary      Get complete anime
// @Description  Scrape and get complete anime from Otakudesu.
// @Produce      json
// @Param        page path int true "Page" Example(2)
// @Success      200 {object} example.GetOdAnimeHomeResponse
// @Router       /otakudesu/complete-anime/page/{page} [get]
func (a *OdAnimeController) GetCompleteAnime(c *fiber.Ctx) error {
	page := c.Params("page")

	completeAnime, err := a.AnimeService.GetCompleteAnime(page)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[model.CompleteAnime]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Successfully Retrieved Anime!",
		Results: completeAnime,
	})
}

// @Tags         Otakudesu
// @Summary      Get ongoing anime
// @Description  Scrape and get ongoing anime from Otakudesu.
// @Produce      json
// @Param        page path int true "Page" Example(2)
// @Success      200 {object} example.GetOdAnimeHomeResponse
// @Router       /otakudesu/ongoing-anime/page/{page} [get]
func (a *OdAnimeController) GetOngoingAnime(c *fiber.Ctx) error {
	page := c.Params("page")

	ongoingAnime, err := a.AnimeService.GetOngoingAnime(page)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[model.OngoingAnime]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Successfully Retrieved Anime!",
		Results: ongoingAnime,
	})
}

// @Tags         Otakudesu
// @Summary      Get trending anime
// @Description  Scrape and get trending anime from Otakudesu.
// @Produce      json
// @Router       /otakudesu/trending [get]
func (a *OdAnimeController) GetTrendingAnime(c *fiber.Ctx) error {
	TrendingAnime, err := a.AnimeService.GetTrendingAnime()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error getting trending anime")
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[model.TrendingAnime]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Successfully Retrieved Trending Anime!",
		Results: TrendingAnime,
	})
}

// @Tags         Otakudesu
// @Summary      Get popular anime
// @Description  Scrape and get popular anime from Otakudesu.
// @Produce      json
// @Router       /otakudesu/popular [get]
func (a *OdAnimeController) GetAnimePopular(c *fiber.Ctx) error {
	PopularAnime, err := a.AnimeService.GetAnimePopular()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error getting popular anime")
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[model.AnimeData]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Successfully Retrieved Popular Anime!",
		Results: PopularAnime,
	})
}

// @Tags         Otakudesu
// @Summary      Get details and episode
// @Description  Scrape and get details and episode from Otakudesu.
// @Produce      json
// @Param        judul path      string  true   "Judul Anime" Example(ds-future-sub-indo)
// @Success      200   {object}  example.GetOdAnimeEpisodeResponse
// @Router       /otakudesu/detail/{judul} [get]
func (a *OdAnimeController) GetAnimeDetails(c *fiber.Ctx) error {
	judul := c.Params("judul")
	detail, episode, err := a.AnimeService.GetAnimeDetails(judul)

	results := model.EpisodePageResult{
		AnimeDetail: detail,
		AnimeEps:    episode,
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.EpisodePageResult]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Success Retrieved Anime!",
		Data:    results,
	})
}

// @Tags         Otakudesu
// @Summary      Get Episode Video Source
// @Description  Scrape and get episode source video from Otakudesu.
// @Produce      json
// @Security BearerAuth
// @Param        judul_eps path string true "Judul Episode" Example(drstn-s4-episode-8-sub-indo)
// @Success      200 {object} example.GetOdAnimeEpisodeVideoResponse
// @Router       /otakudesu/play/{judul_eps} [get]
func (a *OdAnimeController) GetAnimeSourceVid(c *fiber.Ctx) error {
	judul_eps := c.Params("judul_eps")
	animSource, err := a.AnimeService.GetAnimeSourceVid(c, judul_eps)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[model.AnimeSourceData]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Success Retrieved Anime",
		Data:    animSource,
	})
}

// @Tags         Otakudesu
// @Summary      Get Anime Genre
// @Description  Scrape and get anime by genre from Otakudesu.
// @Produce      json
// @Param        genre path string true "Genre Anime" Example(adventure)
// @Param        page path string true "Current Page" Example(0)
// @Success      200 {object} example.GetOdAnimeByGenreResponse
// @Router       /otakudesu/genre/{genre}/page/{page} [get]
func (a *OdAnimeController) GetAnimeGenreList(c *fiber.Ctx) error {
	genre := c.Params("genre")
	page := c.Params("page")

	result, err := a.AnimeService.GetAnimeGenreList(genre, page)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[model.GenreAnime]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Success Retrieved Anime!",
		Results: result,
	})
}

// @Tags         Otakudesu
// @Summary      Get All Anime Genre
// @Description  Scrape and get all anime genre from Otakudesu.
// @Produce      json
// @Router       /otakudesu/genre-list [get]
func (a *OdAnimeController) GetAllGenreList(c *fiber.Ctx) error {
	result, err := a.AnimeService.GetAllAnimeGenre()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "Internal Server Error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[model.GenreList]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Success Retrieved Anime Genre!",
		Results: result,
	})
}

// @Tags         Otakudesu
// @Summary      Search Anime
// @Description  Scrape and search anime by title from Otakudesu.
// @Produce      json
// @Param        title query string true "Title of the Anime" Example(one piece)
// @Success      200 {object} example.GetOdAnimeEpisodeVideoResponse
// @Router       /otakudesu/search [get]
func (a *OdAnimeController) GetAnimeSearchList(c *fiber.Ctx) error {
	title := c.Query("title")

	result, err := a.AnimeService.GetAnimeByTitle(title)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[model.SearchResult]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Success Retrieved Anime!",
		Results: result,
	})
}

package controller

import (
	od_anime_entity "github.com/muhammadsaefulr/NimeStreamAPI/pkg/domain/entity/otakudesu_scrape"
	"github.com/muhammadsaefulr/NimeStreamAPI/pkg/domain/entity/response"

	od_service "github.com/muhammadsaefulr/NimeStreamAPI/pkg/service/otakudesu_scrape"

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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[od_anime_entity.AnimeData]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Successfully Retrieved Anime!",
		Results: animes,
	})
}

// @Tags         Otakudesu
// @Summary      Get details and episode
// @Description  Scrape and get details and episode from Otakudesu.
// @Produce      json
// @Param        judul path      string  true   "Judul Anime" Example(ds-future-sub-indo)
// @Success      200   {object}  example.GetOdAnimeEpisodeResponse
// @Router       /otakudesu/detail/{judul} [get]
func (a *OdAnimeController) GetAnimeEpisode(c *fiber.Ctx) error {
	judul := c.Params("judul")
	detail, episode, err := a.AnimeService.GetAnimeEpisode(judul)

	results := od_anime_entity.EpisodePageResult{
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[od_anime_entity.EpisodePageResult]{
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
// @Param        judul_eps path string true "Judul Episode" Example(drstn-s4-episode-8-sub-indo)
// @Success      200 {object} example.GetOdAnimeEpisodeVideoResponse
// @Router       /otakudesu/play/{judul_eps} [get]
func (a *OdAnimeController) GetAnimeSourceVid(c *fiber.Ctx) error {
	judul_eps := c.Params("judul_eps")
	animSource, err := a.AnimeService.GetAnimeSourceVid(judul_eps)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithDetail[od_anime_entity.AnimeSourceData]{
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[od_anime_entity.GenreAnime]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Success Retrieved Anime!",
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

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithCommonData[od_anime_entity.SearchResult]{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Success Retrieved Anime!",
		Results: result,
	})
}

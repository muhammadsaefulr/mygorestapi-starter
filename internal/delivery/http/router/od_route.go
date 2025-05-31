package router

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/od_controller"
	od_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/otakudesu_scrape"

	"github.com/gofiber/fiber/v2"
)

func OdRoutes(v1 fiber.Router, u od_service.AnimeService) {
	odController := controller.NewAnimeController(u)

	anime := v1.Group("/otakudesu")

	anime.Get("/", odController.GetHomePageAnime)
	anime.Get("/detail/:judul", odController.GetAnimeEpisode)
	anime.Get("/play/:judul_eps", odController.GetAnimeSourceVid)
	anime.Get("/genre/:genre/page/:page", odController.GetAnimeGenreList)
	anime.Get("/search", odController.GetAnimeSearchList)
}

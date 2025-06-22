package router

import (
	controller "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/controller/od_controller"
	m "github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	od_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/otakudesu_scrape"

	"github.com/gofiber/fiber/v2"
)

func OdRoutes(v1 fiber.Router, u od_service.AnimeService) {
	odController := controller.NewAnimeController(u)

	anime := v1.Group("/otakudesu")

	anime.Get("/", odController.GetHomePageAnime)
	anime.Get("/detail/:judul", odController.GetAnimeEpisode)
	anime.Get("/complete-anime/page/:page", odController.GetCompleteAnime)
	anime.Get("/trending", odController.GetTrendingAnime)
	anime.Get("/popular", odController.GetAnimePopular)
	anime.Get("/ongoing-anime/page/:page", odController.GetOngoingAnime)
	anime.Get("/play/:judul_eps", m.Auth(), odController.GetAnimeSourceVid)
	anime.Get("/genre/:genre/page/:page", odController.GetAnimeGenreList)
	anime.Get("/search", odController.GetAnimeSearchList)
}

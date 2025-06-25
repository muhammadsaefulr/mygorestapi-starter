package od_service

import (
	"github.com/gofiber/fiber/v2"
	od_anime_model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type AnimeService interface {
	GetHomePage() ([]od_anime_model.AnimeData, error)
	GetCompleteAnime(page string) ([]od_anime_model.CompleteAnime, error)
	GetOngoingAnime(page string) ([]od_anime_model.OngoingAnime, error)
	GetAnimeDetails(judul string) (od_anime_model.AnimeDetail, []od_anime_model.AnimeEpisode, error)
	GetTrendingAnime() ([]od_anime_model.TrendingAnime, error)
	GetAnimePopular() ([]od_anime_model.AnimeData, error)
	GetAnimeSourceVid(ctx *fiber.Ctx, judul_eps string) (od_anime_model.AnimeSourceData, error)
	GetAnimeGenreList(genre string, page string) ([]od_anime_model.GenreAnime, error)
	GetAllAnimeGenre() ([]od_anime_model.GenreList, error)
	GetAnimeByTitle(title string) ([]od_anime_model.SearchResult, error)
}

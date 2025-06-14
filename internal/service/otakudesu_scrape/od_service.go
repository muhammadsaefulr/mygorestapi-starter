package od_service

import (
	od_anime_model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type AnimeService interface {
	GetHomePage() ([]od_anime_model.AnimeData, error)
	GetAnimeEpisode(judul string) (od_anime_model.AnimeDetail, []od_anime_model.AnimeEpisode, error)
	GetAnimeSourceVid(judul_eps string) (od_anime_model.AnimeSourceData, error)
	GetAnimeGenreList(genre string, page string) ([]od_anime_model.GenreAnime, error)
	GetAnimeByTitle(title string) ([]od_anime_model.SearchResult, error)
}

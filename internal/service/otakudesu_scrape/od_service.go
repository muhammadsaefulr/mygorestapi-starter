package od_service

import (
	od_anime_entity "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/entity/otakudesu_scrape"
)

type AnimeService interface {
	GetHomePage() ([]od_anime_entity.AnimeData, error)
	GetAnimeEpisode(judul string) (od_anime_entity.AnimeDetail, []od_anime_entity.AnimeEpisode, error)
	GetAnimeSourceVid(judul_eps string) (od_anime_entity.AnimeSourceData, error)
	GetAnimeGenreList(genre string, page string) ([]od_anime_entity.GenreAnime, error)
	GetAnimeByTitle(title string) ([]od_anime_entity.SearchResult, error)
}

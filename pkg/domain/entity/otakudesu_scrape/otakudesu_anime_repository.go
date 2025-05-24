package od_anime_entity

type AnimeRepository interface {
	ScrapeHomePage() ([]AnimeData, error)
	ScrapeAnimepisodes(judul string) ([]AnimeEpisode, error)
	ScrapeAnimeSourceData(judul_eps string) (AnimeSourceData, error)
}

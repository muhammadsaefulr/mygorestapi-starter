package od_service

import (
	od_anime_entity "github.com/muhammadsaefulr/NimeStreamAPI/pkg/domain/entity/otakudesu_scrape"
	modules "github.com/muhammadsaefulr/NimeStreamAPI/pkg/infrastructure/modules/scrape_otakudesu"
)

type AnimeService interface {
	GetHomePage() ([]od_anime_entity.AnimeData, error)
	GetAnimeEpisode(judul string) (od_anime_entity.AnimeDetail, []od_anime_entity.AnimeEpisode, error)
	GetAnimeSourceVid(judul_eps string) (od_anime_entity.AnimeSourceData, error)
	GetAnimeGenreList(genre string, page string) ([]od_anime_entity.GenreAnime, error)
	GetAnimeByTitle(title string) ([]od_anime_entity.SearchResult, error)
}

type animeService struct{}

func NewAnimeService() AnimeService {
	return &animeService{}
}

var mainUrl = "https://otakudesu.cloud"

func (s *animeService) GetHomePage() ([]od_anime_entity.AnimeData, error) {
	animes := modules.ScrapeHomePage()

	return animes, nil
}

func (s *animeService) GetAnimeEpisode(judul string) (od_anime_entity.AnimeDetail, []od_anime_entity.AnimeEpisode, error) {
	detail, eps := modules.ScrapeAnimeEpisodes(mainUrl + ("/anime/" + judul))

	return detail, eps, nil
}

func (s *animeService) GetAnimeSourceVid(judul_eps string) (od_anime_entity.AnimeSourceData, error) {
	animSource := modules.ScrapeAnimeSourceData(mainUrl + ("/episode/" + judul_eps))

	return animSource, nil
}

func (s *animeService) GetAnimeGenreList(genre string, page string) ([]od_anime_entity.GenreAnime, error) {
	animGenre := modules.ScrapeGenreAnime(mainUrl + "/genres/" + (genre + "/page/" + page))

	return animGenre, nil
}

func (s *animeService) GetAnimeByTitle(title string) ([]od_anime_entity.SearchResult, error) {
	animSearch := modules.ScrapeSearchAnimeByTitle(mainUrl + "?s=" + title + "&post_type=anime")

	return animSearch, nil
}

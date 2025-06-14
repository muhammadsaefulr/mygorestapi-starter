package od_service

import (
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	modules "github.com/muhammadsaefulr/NimeStreamAPI/internal/infrastructure/modules/scrape_otakudesu"
)

type animeService struct{}

func NewAnimeService() AnimeService {
	return &animeService{}
}

var mainUrl = "https://otakudesu.cloud"

func (s *animeService) GetHomePage() ([]model.AnimeData, error) {
	animes := modules.ScrapeHomePage()

	return animes, nil
}

func (s *animeService) GetAnimeEpisode(judul string) (model.AnimeDetail, []model.AnimeEpisode, error) {
	detail, eps := modules.ScrapeAnimeEpisodes(mainUrl + ("/anime/" + judul))

	return detail, eps, nil
}

func (s *animeService) GetAnimeSourceVid(judul_eps string) (model.AnimeSourceData, error) {
	animSource := modules.ScrapeAnimeSourceData(mainUrl + ("/episode/" + judul_eps))

	return animSource, nil
}

func (s *animeService) GetAnimeGenreList(genre string, page string) ([]model.GenreAnime, error) {
	animGenre := modules.ScrapeGenreAnime(mainUrl + "/genres/" + (genre + "/page/" + page))

	return animGenre, nil
}

func (s *animeService) GetAnimeByTitle(title string) ([]model.SearchResult, error) {
	animSearch := modules.ScrapeSearchAnimeByTitle(mainUrl + "?s=" + title + "&post_type=anime")

	return animSearch, nil
}

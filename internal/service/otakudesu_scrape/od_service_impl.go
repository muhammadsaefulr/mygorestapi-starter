package od_service

import (
	"sort"
	"strconv"
	"strings"

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

func (s *animeService) GetCompleteAnime(page string) ([]model.CompleteAnime, error) {
	ongoinAnime := modules.ScrapeCompleteAnime(page)

	return ongoinAnime, nil
}

func (s *animeService) GetOngoingAnime(page string) ([]model.OngoingAnime, error) {
	ongoingAnime := modules.ScrapeOngoingAnime(page)

	return ongoingAnime, nil
}

func (s *animeService) GetTrendingAnime() ([]model.TrendingAnime, error) {
	var results []model.TrendingAnime

	ongoing := modules.ScrapeOngoingAnime("1")
	for _, o := range ongoing {
		if strings.Contains(o.DaysUpdated, "Rabu") ||
			strings.Contains(o.DaysUpdated, "Minggu") {
			results = append(results, model.TrendingAnime{
				Title:        o.Title,
				URL:          o.URL,
				JudulPath:    o.JudulPath,
				ThumbnailURL: o.ThumbnailURL,
				LatestEp:     o.Episode,
				UpdateAnime:  o.UpdatedAt,
			})
		}
	}

	complete := modules.ScrapeCompleteAnime("1")
	for _, c := range complete {
		rating, err := strconv.ParseFloat(c.Rating, 64)
		if err != nil {
			rating = 0.0
		}
		if rating >= 7.5 {
			results = append(results, model.TrendingAnime{
				Title:        c.Title,
				URL:          c.URL,
				JudulPath:    c.JudulPath,
				ThumbnailURL: c.ThumbnailURL,
				LatestEp:     c.LatestEp,
				UpdateAnime:  c.UpdatedAt,
			})
		}
	}

	// Sort berdasarkan tanggal update terbaru (descending)
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].UpdateAnime > results[j].UpdateAnime
	})

	return results, nil
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

package od_service

import (
	"context"
	"fmt"
	"net/url"

	// "log"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	modules "github.com/muhammadsaefulr/NimeStreamAPI/internal/infrastructure/modules/scrape_otakudesu"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/track_episode_view"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
)

type animeService struct {
	TrackEpisodeViewRepo repository.TrackEpisodeViewRepository
}

func NewAnimeService(trackRepo repository.TrackEpisodeViewRepository) AnimeService {
	return &animeService{
		TrackEpisodeViewRepo: trackRepo,
	}
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

func (s *animeService) GetAnimePopular() ([]model.AnimeData, error) {
	var animeData []model.AnimeData
	ctx := context.Background()

	// --- PRIORITASKAN DUMMY SCRAPE ---
	dummySlugs := []string{
		"1piece-sub-indo",
		"bleach-oukoku-tan-sub-indo",
		"kimetsu-yaiba-subtitle-indonesia",
		"one-punch-sub-indo",
		"dea-note-subtitle-indonesia",
		"fulltal-alchemist-sub-indo",
	}

	var wgDummy sync.WaitGroup
	dummyChan := make(chan model.AnimeData, len(dummySlugs))

	for _, slug := range dummySlugs {
		wgDummy.Add(1)
		go func(slug string) {
			defer wgDummy.Done()
			detail, _, _ := modules.ScrapeAnimeDetail(mainUrl + "/anime/" + slug)
			dummyChan <- model.AnimeData{
				Title:        detail.Title,
				URL:          slug,
				ThumbnailURL: detail.ThumbnailURL,
				LatestEp:     detail.TotalEps,
				UpdateAnime:  detail.ReleaseDate,
				JudulPath:    slug,
			}
		}(slug)
	}

	go func() {
		wgDummy.Wait()
		close(dummyChan)
	}()

	for data := range dummyChan {
		animeData = append(animeData, data)
	}

	// --- SCRAPE DARI TRACK VIEW (POPULAR) ---
	topEpisodes, err := s.TrackEpisodeViewRepo.GetAll(ctx, &repository.TrackEpsParam{
		Page:      1,
		Limit:     10,
		MovieType: "anime",
	})

	if err != nil {
		return animeData, nil
	}

	// --- Worker Pool Configuration ---
	const maxWorkers = 10
	jobs := make(chan repository.TrackEpisodeViewSummary, len(topEpisodes))
	results := make(chan model.AnimeData, len(topEpisodes))

	// Start worker goroutines
	for w := 0; w < maxWorkers; w++ {
		go func() {
			for job := range jobs {
				slug := path.Base(strings.TrimSuffix(job.MovieDetailUrl, "/"))
				detail, _, _ := modules.ScrapeAnimeDetail(mainUrl + "/anime/" + slug)

				results <- model.AnimeData{
					Title:        detail.Title,
					URL:          job.EpisodeId,
					ThumbnailURL: detail.ThumbnailURL,
					LatestEp:     detail.TotalEps,
					UpdateAnime:  detail.ReleaseDate,
					JudulPath:    slug,
				}
			}
		}()
	}

	// Kirim semua pekerjaan ke channel
	go func() {
		for _, top := range topEpisodes {
			jobs <- top
		}
		close(jobs)
	}()

	// Terima semua hasil
	for i := 0; i < len(topEpisodes); i++ {
		animeData = append(animeData, <-results)
	}

	if len(animeData) > 15 {
		animeData = animeData[:15]
	}

	return animeData, nil
}

func (s *animeService) GetAnimeDetails(judul string) (model.AnimeDetail, []model.AnimeEpisode, []model.SearchResult, error) {

	detail, eps, rekomend := modules.ScrapeAnimeDetail(mainUrl + ("/anime/" + judul))

	return detail, eps, rekomend, nil
}

func (s *animeService) GetAnimeSourceVid(ctx *fiber.Ctx, judul_eps string) (model.AnimeSourceData, error) {
	authHeader := ctx.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	IdUsr, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return model.AnimeSourceData{}, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("Error verifying token: %s", err.Error()))
	}

	userUUID, err := uuid.Parse(IdUsr)
	if err != nil {
		return model.AnimeSourceData{}, fiber.NewError(fiber.StatusBadRequest, "Invalid user ID format")
	}

	animSource := modules.ScrapeAnimeSourceData(mainUrl + ("/episode/" + judul_eps))

	s.TrackEpisodeViewRepo.Create(context.Background(), model.TrackEpisodeView{
		UserId:         userUUID,
		TypeMovie:      "anime",
		MovieDetailUrl: animSource.DetailURL,
		EpisodeId:      judul_eps,
	})

	return animSource, nil
}

func (s *animeService) GetAnimeGenreList(genre string, page string) ([]model.GenreAnime, error) {
	animGenre := modules.ScrapeGenreAnime(mainUrl + "/genres/" + (genre + "/page/" + page))

	return animGenre, nil
}

func (s *animeService) GetAllAnimeGenre() ([]model.GenreList, error) {
	animGenre := modules.ScrapeGenreList(mainUrl + "/genre-list")

	return animGenre, nil
}

func (s *animeService) GetAnimeByTitle(title string) ([]model.SearchResult, error) {
	animSearch := modules.ScrapeSearchAnimeByTitle(mainUrl + "?s=" + url.QueryEscape(title) + "&post_type=anime")

	return animSearch, nil
}

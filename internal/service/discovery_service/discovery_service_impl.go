package service

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/discovery/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	repository "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/discovery"
	svcAnilist "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/anilist_service"
	svcMdl "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/mdl_service"
	svcMovieDt "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/movie_details_service"
	od_service "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/otakudesu_scrape"
	svcTmdb "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/tmdb_service"

	// "github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type DiscoveryService struct {
	Log        *logrus.Logger
	Validate   *validator.Validate
	Repo       repository.DiscoveryRepo
	AnilistSvc svcAnilist.AnilistServiceInterface
	TmdbSvc    svcTmdb.TmdbServiceInterface
	MdlSvc     svcMdl.MdlServiceInterface
	OdService  od_service.AnimeService
	MovieSvc   svcMovieDt.MovieDetailsServiceInterface
	redisCl    *redis.Client
}

// GetDiscoverSearch implements DiscoveryServiceInterface.

func NewDiscoveryService(validate *validator.Validate, redisCl *redis.Client, SvcAn svcAnilist.AnilistServiceInterface, svcTmdb svcTmdb.TmdbServiceInterface, svcMdl svcMdl.MdlServiceInterface, svcOd od_service.AnimeService, svcMvDt svcMovieDt.MovieDetailsServiceInterface) DiscoveryServiceInterface {
	return &DiscoveryService{
		Log:        utils.Log,
		Validate:   validate,
		AnilistSvc: SvcAn,
		TmdbSvc:    svcTmdb,
		MdlSvc:     svcMdl,
		OdService:  svcOd,
		MovieSvc:   svcMvDt,
		redisCl:    redisCl,
	}
}

func (s *DiscoveryService) GetDiscover(c *fiber.Ctx, params *request.QueryDiscovery) ([]response.MovieDetailOnlyResponse, int64, int64, error) {
	log.Printf("param info: %+v", params)

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid query params")
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	key := BuildRedisKeyFromDiscoveryParams(params)

	type cacheResult struct {
		Data      []response.MovieDetailOnlyResponse `json:"data"`
		PageTotal int64                              `json:"page"`
		TotalData int64                              `json:"total"`
	}

	var cached cacheResult
	if err := s.GetDiscoveryCache(c.Context(), key, &cached); err == nil {
		return cached.Data, cached.PageTotal, cached.TotalData, nil
	}

	var (
		results  = []response.MovieDetailOnlyResponse{}
		pageRes  int64
		TotalRes int
		firstErr error
		mu       sync.Mutex
		wg       sync.WaitGroup
	)

	if strings.ToLower(params.Category) == "search" {
		log.Println(params.Type)
		switch strings.ToLower(params.Type) {
		case "anime":
			wg.Add(1)
			go func() {
				defer wg.Done()
				tmp := *params

				movieData, _, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(&tmp))
				if err != nil {
					s.Log.Errorf("MovieSvc.GetAll error (anime): %v", err)
				}

				anilistResults, _, err := s.AnilistSvc.GetAll(c, convert_types.MapToAnilistQuery(&tmp))
				if err != nil || len(anilistResults) == 0 {
					s.Log.Errorf("Anilist fetch error: %v", err)
					return
				}

				odResults, err := s.OdService.GetAnimeByTitle(tmp.Search)
				if err != nil || len(odResults) == 0 {
					s.Log.Errorf("OD fetch error: %v", err)
				}

				var aniTitles []string
				for _, a := range anilistResults {
					aniTitles = append(aniTitles, a.Title)
				}

				var odTitles []string
				for _, o := range odResults {
					odTitles = append(odTitles, o.Title)
				}

				aIdx, bIdx := utils.JaroWinklerPairIndices(aniTitles, odTitles, 5)

				var picked []response.MovieDetailOnlyResponse
				if len(aIdx) > 0 {
					for i := range aIdx {
						ani := anilistResults[aIdx[i]]
						if len(movieData) > 0 {
							ani.MovieID = movieData[0].MovieID
							ani.PathURL = "/movie/details/" + movieData[0].MovieID
						} else if len(bIdx) > i && bIdx[i] < len(odResults) {
							od := odResults[bIdx[i]]
							ani.MovieID = path.Base(strings.TrimSuffix(od.URL, "/"))
							ani.PathURL = od.URL
							ani.ThumbnailURL = od.ThumbnailURL
							ani.Status = od.Status
							ani.Rating = od.Rating
							ani.Genres = []string{}
							for _, g := range od.Genres {
								ani.Genres = append(ani.Genres, g.Title)
							}
						}
						picked = append(picked, ani)
					}
				} else {
					// fallback ke anilistResults[0]
					main := anilistResults[0]
					if len(movieData) > 0 {
						main.MovieID = movieData[0].MovieID
						main.PathURL = "/movie/details/" + movieData[0].MovieID
					}
					picked = append(picked, main)
				}

				if len(picked) == 0 || picked[0].PathURL == "" {
					mu.Lock()
					firstErr = fiber.NewError(fiber.StatusNotFound, "Tidak ditemukan di MovieDetails maupun OD")
					mu.Unlock()
					return
				}

				mu.Lock()
				results = append(results, picked...)
				pageRes = 1
				TotalRes = len(results)
				mu.Unlock()
			}()
		case "movie", "tv":
			wg.Add(1)
			go func() {
				defer wg.Done()
				tmp := *params

				movieData, _, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(&tmp))
				if err != nil {
					s.Log.Errorf("MovieSvc.GetAll error (%s): %v", tmp.Type, err)
				}

				tmdbResults, _, _, err := s.TmdbSvc.GetAll(c, convert_types.MapToTmdbQuery(&tmp))
				if err != nil || len(tmdbResults) == 0 {
					s.Log.Errorf("TMDb fetch error: %v", err)
					return
				}

				var tmdbTitles []string
				for _, t := range tmdbResults {
					tmdbTitles = append(tmdbTitles, t.Title)
				}

				var dbTitles []string
				for _, m := range movieData {
					dbTitles = append(dbTitles, m.Title)
				}

				var picked []response.MovieDetailOnlyResponse
				if len(dbTitles) > 0 {
					aIdx, bIdx := utils.JaroWinklerPairIndices(tmdbTitles, dbTitles, 5)

					for i := range aIdx {
						main := tmdbResults[aIdx[i]]
						main.MovieID = movieData[bIdx[i]].MovieID
						main.PathURL = "/movie/details/" + movieData[bIdx[i]].MovieID
						picked = append(picked, main)
					}
				}

				if len(picked) == 0 {
					main := tmdbResults[0]
					if len(movieData) > 0 {
						main.MovieID = movieData[0].MovieID
						main.PathURL = "/movie/details/" + movieData[0].MovieID
					}
					picked = append(picked, main)
				}

				if len(picked) == 0 || picked[0].PathURL == "" {
					mu.Lock()
					firstErr = fiber.NewError(fiber.StatusNotFound, "Tidak ditemukan di MovieDetails")
					mu.Unlock()
					return
				}

				mu.Lock()
				results = append(results, picked...)
				pageRes = 1
				TotalRes = len(results)
				mu.Unlock()
			}()

		}

		wg.Wait()
		if firstErr != nil && len(results) == 0 {
			return nil, 0, 0, firstErr
		}

		return results, pageRes, int64(TotalRes), nil
	}

	switch strings.ToLower(params.Type) {
	case "anime":
		wg.Add(1)
		go func() {
			defer wg.Done()

			data, totalPage, err := s.AnilistSvc.GetAll(c, convert_types.MapToAnilistQuery(params))
			if err != nil {
				s.Log.Errorf("AnilistSvc.GetAll error: %v", err)
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			var (
				filtered   = make([]response.MovieDetailOnlyResponse, 0, len(data))
				numRegex   = regexp.MustCompile(`\d+`)
				limitCount = 0
			)

			for _, d := range data {
				if limitCount >= params.Limit {
					break
				}

				tmp := *params
				tmp.Search = d.Title
				movieData, _, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(&tmp))
				if err == nil && len(movieData) > 0 {
					d.MovieID = movieData[0].MovieID
					d.PathURL = "/movie/details/" + movieData[0].MovieID
					filtered = append(filtered, d)
					limitCount++
					continue
				}

				odResults, err := s.OdService.GetAnimeByTitle(d.Title)
				if err != nil || len(odResults) == 0 {
					continue
				}

				dTitleLower := strings.ToLower(d.Title)
				dTitleWords := strings.Fields(dTitleLower)

				if len(dTitleWords) > 4 {
					dTitleWords = dTitleWords[:4]
				}
				matchKey := strings.Join(dTitleWords, " ")
				dNums := numRegex.FindAllString(dTitleLower, -1)

				for _, od := range odResults {
					if od.URL == "" {
						continue
					}

					odTitleLower := strings.ToLower(od.Title)
					if !strings.Contains(odTitleLower, matchKey) {
						continue
					}

					numMatch := true
					for _, num := range dNums {
						if !strings.Contains(odTitleLower, num) {
							numMatch = false
							break
						}
					}

					if len(dNums) == 0 || numMatch {
						d.MovieID = path.Base(strings.TrimSuffix(od.URL, "/"))
						d.PathURL = od.URL
						d.Title = od.Title
						d.ThumbnailURL = od.ThumbnailURL
						filtered = append(filtered, d)
						limitCount++
						break
					}
				}
			}

			mu.Lock()
			results = append(results, filtered...)
			TotalRes = len(filtered)
			pageRes = totalPage
			mu.Unlock()
		}()

	case "movie", "tv":
		wg.Add(1)
		go func() {
			defer wg.Done()

			data, _, totalPages, err := s.TmdbSvc.GetAll(c, convert_types.MapToTmdbQuery(params))
			if err != nil {
				s.Log.Errorf("TmdbSvc.GetAll error: %v", err)
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			var filtered []response.MovieDetailOnlyResponse
			for _, d := range data {
				tmp := *params
				tmp.Search = d.Title
				movieData, _, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(&tmp))
				if err == nil && len(movieData) > 0 {
					d.MovieID = movieData[0].MovieID
					filtered = append(filtered, d)
				}
			}

			mu.Lock()
			results = append(results, filtered...)
			pageRes = totalPages
			TotalRes = len(results)

			mu.Unlock()
		}()

	case "kdrama":
		wg.Add(1)
		go func() {
			defer wg.Done()

			data, page, _, err := s.MdlSvc.GetAll(c, convert_types.MapToMdlQuery(params))
			if err != nil {
				s.Log.Errorf("MdlSvc.GetAll error: %v", err)
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			var filtered []response.MovieDetailOnlyResponse
			for _, d := range data {
				tmp := *params
				tmp.Search = d.Title
				movieData, _, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(&tmp))
				if err == nil && len(movieData) > 0 {
					d.MovieID = movieData[0].MovieID
					d.PathURL = movieData[0].PathURL
					filtered = append(filtered, d)
				}
			}

			mu.Lock()
			results = append(results, filtered...)
			pageRes = page
			TotalRes = len(filtered)
			mu.Unlock()
		}()

	default:
		return nil, 0, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid query params")
	}

	wg.Wait()

	if firstErr != nil && len(results) == 0 {
		return nil, 0, 0, firstErr
	}

	_ = s.SetDiscoveryCache(c.Context(), key, cacheResult{
		Data:      results,
		PageTotal: pageRes,
		TotalData: int64(TotalRes),
	}, params)

	return results, pageRes, int64(TotalRes), nil
}

func (s *DiscoveryService) GetDiscoverDetailByTitle(c *fiber.Ctx, mediaType string, slug string) (*response.MovieDetailOnlyResponse, error) {
	title := url.QueryEscape(strings.ReplaceAll(slug, "-", " "))
	log.Println(mediaType)

	params := &request.QueryDiscovery{
		Search:   title,
		Category: "search",
		Type:     mediaType,
		Page:     1,
		Limit:    1,
	}

	results, _, _, err := s.GetDiscover(c, params)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Failed to search discovery: "+err.Error())
	}
	if len(results) == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "No result found for title: "+title)
	}

	first := results[0]
	id := first.IDSource
	media := first.MovieType

	var detail *response.MovieDetailOnlyResponse

	switch strings.ToLower(mediaType) {
	case "anime":
		detail, err = s.AnilistSvc.GetMovieDetailsByID(c, id)
	case "movie", "tv":
		detail, err = s.TmdbSvc.GetDetailByID(c, id, media)
	case "kdrama":
		detail, err = s.MdlSvc.GetDetailByID(c, id)
	default:
		return nil, fiber.NewError(fiber.StatusBadRequest, "Unsupported mediaType: "+media)
	}

	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadGateway, "Failed to fetch detail: "+err.Error())
	}

	if detail != nil {
		searchParam := &request.QueryDiscovery{Search: detail.Title}
		movieData, _, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(searchParam))
		if err == nil && len(movieData) > 0 {
			detail.MovieID = movieData[0].MovieID
			detail.PathURL = "/movie/details/" + detail.MovieID
		} else if media == "anime" {
			odData, err := s.OdService.GetAnimeByTitle(detail.Title)
			if err == nil && len(odData) > 0 {
				detail.MovieID = path.Base(strings.TrimSuffix(odData[0].URL, "/"))
				detail.PathURL = "/otakudesu/detail/" + detail.MovieID
			}
		}
	}

	if detail.Rekomend != nil {
		var finalRekom []response.MovieDetailOnlyResponse

		for _, r := range *detail.Rekomend {
			found := false

			searchParam := &request.QueryDiscovery{Search: r.Title}
			movieData, _, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(searchParam))
			if err == nil && len(movieData) > 0 {
				r.MovieID = r.Title
				r.PathURL = "/movie/details/" + r.MovieID
				found = true
			} else if media == "anime" {
				odData, err := s.OdService.GetAnimeByTitle(r.Title)
				if err == nil && len(odData) > 0 {
					r.MovieID = odData[0].URL
					r.PathURL = "/otakudesu/detail/" + r.MovieID
					found = true
				}
			}

			if found {
				finalRekom = append(finalRekom, r)
			}
		}

		if len(finalRekom) > 0 {
			detail.Rekomend = &finalRekom
		} else {
			detail.Rekomend = nil
		}
	}

	return detail, nil
}

func (s *DiscoveryService) GetDiscoverGenres(c *fiber.Ctx, params *request.QueryDiscovery) ([]response.GenreDetail, error) {

	if err := s.Validate.Struct(params); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid query params")
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	key := BuildRedisKeyFromDiscoveryParams(params)

	type cacheResult struct {
		Data      []response.GenreDetail `json:"data"`
		PageTotal int64                  `json:"page"`
		TotalData int64                  `json:"total"`
	}

	var cached cacheResult
	if err := s.GetDiscoveryCache(c.Context(), key, &cached); err == nil {
		return cached.Data, nil
	}

	var (
		results  = []response.GenreDetail{}
		firstErr error
		mu       sync.Mutex
		wg       sync.WaitGroup
	)

	switch strings.ToLower(params.Type) {
	case "anime":
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := s.AnilistSvc.GetAllGenres(c)
			if err != nil {
				s.Log.Errorf("AnilistSvc.GetAllGenres error: %v", err)
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}
			var filtered []response.GenreDetail
			for _, d := range data {
				filtered = append(filtered, d)
			}
			mu.Lock()
			results = append(results, filtered...)
			mu.Unlock()
		}()
	case "movie":
		wg.Add(1)
		go func() {
			defer wg.Done()

			data, err := s.TmdbSvc.GetAllGenres(c)
			if err != nil {
				s.Log.Errorf("TmdbSvc.GetAllGenres error: %v", err)
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			var filtered []response.GenreDetail
			for _, d := range data {
				filtered = append(filtered, d)
			}

			mu.Lock()
			results = append(results, filtered...)
			mu.Unlock()
		}()

	case "kdrama":
		wg.Add(1)
		go func() {
			defer wg.Done()

			data, err := s.MdlSvc.GetAllGenres(c)
			if err != nil {
				s.Log.Errorf("MdlSvc.GetAllGenres error: %v", err)
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}
			var filtered []response.GenreDetail
			for _, d := range data {
				filtered = append(filtered, d)
			}

			mu.Lock()
			results = append(results, filtered...)
			mu.Unlock()
		}()

	default:
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid query params")
	}

	wg.Wait()

	if firstErr != nil && len(results) == 0 {
		return nil, firstErr
	}

	_ = s.SetDiscoveryCache(c.Context(), key, cacheResult{
		Data:      results,
		PageTotal: 1,
		TotalData: int64(len(results)),
	}, params)

	return results, nil
}

// Utils

func BuildRedisKeyFromDiscoveryParams(params *request.QueryDiscovery) string {
	parts := []string{"discovery"}

	if params.Type != "" {
		parts = append(parts, "type", params.Type)
	}
	if params.Category != "" {
		parts = append(parts, "category", params.Category)
	}
	if params.Genre != "" {
		parts = append(parts, "genre", params.Genre)
	}
	if params.Search != "" {
		parts = append(parts, "search", params.Search)
	}
	if params.Page > 0 {
		parts = append(parts, "page", strconv.Itoa(params.Page))
	}
	if params.Limit > 0 {
		parts = append(parts, "limit", strconv.Itoa(params.Limit))
	}

	return strings.Join(parts, ":")
}

func (s *DiscoveryService) GetDiscoveryCache(ctx context.Context, key string, dst any) error {
	data, err := s.redisCl.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dst)
}

func (s *DiscoveryService) SetDiscoveryCache(ctx context.Context, key string, val any, params *request.QueryDiscovery) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	ttl := time.Hour
	if params.Search != "" || params.Genre != "" {
		ttl = 30 * time.Minute
	} else if params.Category == "popular" {
		ttl = 6 * time.Hour
	} else if params.Category == "ongoing" || params.Category == "trending" {
		ttl = 1 * time.Hour
	}

	return s.redisCl.Set(ctx, key, b, ttl).Err()
}

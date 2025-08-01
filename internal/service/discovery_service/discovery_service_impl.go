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
	"github.com/xrash/smetrics"

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
				if err != nil {
					s.Log.Errorf("Anilist fetch error: %v", err)
				}

				if len(movieData) > 0 {
					base := movieData[0]

					bestIdx := -1
					bestScore := -1.0
					norm := func(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
					score := func(a, b string) float64 {
						a, b = norm(a), norm(b)
						if a == "" || b == "" {
							return 0
						}
						s := smetrics.JaroWinkler(a, b, 0.7, 4)
						if strings.Contains(a, b) || strings.Contains(b, a) {
							s += 0.05
						}
						return s
					}
					for i := range anilistResults {
						sc := score(base.Title, anilistResults[i].Title)
						if sc > bestScore {
							bestScore = sc
							bestIdx = i
						}
					}
					if bestIdx != -1 {
						ani := anilistResults[bestIdx]
						if base.ThumbnailURL == "" && ani.ThumbnailURL != "" {
							base.ThumbnailURL = ani.ThumbnailURL
						}
						if base.Synopsis == "" && ani.Synopsis != "" {
							base.Synopsis = ani.Synopsis
						}
						if len(base.Genres) == 0 && len(ani.Genres) > 0 {
							base.Genres = ani.Genres
						}
						if base.Rating == "" && ani.Rating != "" {
							base.Rating = ani.Rating
						}
						if base.Status == "" && ani.Status != "" {
							base.Status = ani.Status
						}
					}

					if base.PathURL == "" {
						base.PathURL = "/movie/details/" + base.MovieID
					}
					if base.PathURL == "" {
						mu.Lock()
						firstErr = fiber.NewError(fiber.StatusNotFound, "Tidak ditemukan di MovieDetails")
						mu.Unlock()
						return
					}

					mu.Lock()
					results = append(results, base)
					pageRes = 1
					TotalRes = len(results)
					mu.Unlock()
					return
				}

				odResults, err := s.OdService.GetAnimeByTitle(tmp.Search)
				if err != nil || len(odResults) == 0 {
					mu.Lock()
					firstErr = fiber.NewError(fiber.StatusNotFound, "Tidak ditemukan di MovieDetails maupun Otakudesu")
					mu.Unlock()
					return
				}

				var odTitles, aniTitles []string
				for _, o := range odResults {
					odTitles = append(odTitles, o.Title)
				}
				for _, a := range anilistResults {
					aniTitles = append(aniTitles, a.Title)
				}

				aniIdxPerOD := utils.MatchSourceIndices(odTitles, aniTitles, 0.60)

				var picked []response.MovieDetailOnlyResponse
				for oi, od := range odResults {
					var out response.MovieDetailOnlyResponse

					out.Title = od.Title
					out.PathURL = od.URL
					out.MovieType = "anime"
					out.MovieID = path.Base(strings.TrimSuffix(od.URL, "/"))
					out.ThumbnailURL = od.ThumbnailURL
					out.Status = od.Status
					out.Rating = od.Rating

					for _, g := range od.Genres {
						out.Genres = append(out.Genres, g.Title)
					}

					ai := aniIdxPerOD[oi]
					if ai >= 0 && ai < len(anilistResults) {
						ani := anilistResults[ai]
						if out.ThumbnailURL == "" && ani.ThumbnailURL != "" {
							out.ThumbnailURL = ani.ThumbnailURL
						}
						if out.Synopsis == "" && ani.Synopsis != "" {
							out.Synopsis = ani.Synopsis
						}
						if len(out.Genres) == 0 && len(ani.Genres) > 0 {
							out.Genres = ani.Genres
						}
						if out.Rating == "" && ani.Rating != "" {
							out.Rating = ani.Rating
						}
						if out.Status == "" && ani.Status != "" {
							out.Status = ani.Status
						}
						if out.TotalEps == "" && ani.TotalEps != "" {
							out.TotalEps = ani.TotalEps
						}
						if out.Studio == "" && ani.Studio != "" {
							out.Studio = ani.Studio
						}
						if out.ReleaseDate == "" && ani.ReleaseDate != "" {
							out.ReleaseDate = utils.ConvertDateStripToDay(ani.ReleaseDate)
						}
					}

					picked = append(picked, out)
				}

				if len(picked) == 0 || picked[0].PathURL == "" {
					mu.Lock()
					firstErr = fiber.NewError(fiber.StatusNotFound, "Tidak ditemukan di MovieDetails maupun Otakudesu")
					mu.Unlock()
					return
				}

				mu.Lock()
				results = append(results, picked...)
				pageRes = 1
				TotalRes = len(results)
				mu.Unlock()
			}()

		case "movie", "kdrama":
			wg.Add(1)
			go func() {
				defer wg.Done()
				tmp := *params

				movieData, _, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(&tmp))
				if err != nil {
					s.Log.Errorf("MovieSvc.GetAll error (%s): %v", tmp.Type, err)
					mu.Lock()
					firstErr = fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data Movie Details")
					mu.Unlock()
					return
				}

				if len(movieData) == 0 {
					mu.Lock()
					firstErr = fiber.NewError(fiber.StatusNotFound, "MovieDetails tidak ditemukan")
					mu.Unlock()
					return
				}

				var picked []response.MovieDetailOnlyResponse
				for _, m := range movieData {
					picked = append(picked, response.MovieDetailOnlyResponse{
						MovieID:      m.MovieID,
						Title:        m.Title,
						PathURL:      "/movie/details/" + m.MovieID,
						ThumbnailURL: m.ThumbnailURL,
						Rating:       m.Rating,
						ReleaseDate:  utils.ConvertDateStripToDay(m.ReleaseDate),
						Genres:       m.Genres,
						MovieType:    m.MovieType,
						Synopsis:     m.Synopsis,
						Status:       m.Status,
						Producer:     m.Producer,
						Studio:       m.Studio,
						TotalEps:     m.TotalEps,
					})
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

	case "movie":
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
					d.PathURL = "/movie/details/" + movieData[0].MovieID
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
					d.PathURL = "/movie/details/" + movieData[0].MovieID
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

			filtered = append(filtered, data...)

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

			filtered = append(filtered, data...)

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

			filtered = append(filtered, data...)

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

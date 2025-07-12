package service

import (
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

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
}

// GetDiscoverSearch implements DiscoveryServiceInterface.

func NewDiscoveryService(validate *validator.Validate, SvcAn svcAnilist.AnilistServiceInterface, svcTmdb svcTmdb.TmdbServiceInterface, svcMdl svcMdl.MdlServiceInterface, svcOd od_service.AnimeService, svcMvDt svcMovieDt.MovieDetailsServiceInterface) DiscoveryServiceInterface {
	return &DiscoveryService{
		Log:        utils.Log,
		Validate:   validate,
		AnilistSvc: SvcAn,
		TmdbSvc:    svcTmdb,
		MdlSvc:     svcMdl,
		OdService:  svcOd,
		MovieSvc:   svcMvDt,
	}
}

func (s *DiscoveryService) GetDiscover(c *fiber.Ctx, params *request.QueryDiscovery) ([]response.MovieDetailOnlyResponse, int64, error) {
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid query params")
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	var (
		results  = []response.MovieDetailOnlyResponse{}
		pageRes  int64
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
				main := anilistResults[0]

				if len(movieData) > 0 {
					main.MovieID = movieData[0].MovieID
					main.PathURL = "/movie/details/" + movieData[0].MovieID
				} else {
					odResults, err := s.OdService.GetAnimeByTitle(tmp.Search)
					if err == nil && len(odResults) > 0 {
						main.MovieID = odResults[0].URL
						main.PathURL = "/otakudesu/detail/" + odResults[0].URL
					}
				}

				if main.PathURL == "" {
					mu.Lock()
					firstErr = fiber.NewError(fiber.StatusNotFound, "Tidak ditemukan di MovieDetails maupun OD")
					mu.Unlock()
					return
				}

				mu.Lock()
				results = append(results, main)
				pageRes = 1
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

				tmdbResults, _, err := s.TmdbSvc.GetAll(c, convert_types.MapToTmdbQuery(&tmp))
				if err != nil || len(tmdbResults) == 0 {
					s.Log.Errorf("TMDb fetch error: %v", err)
					return
				}
				main := tmdbResults[0]

				if len(movieData) > 0 {
					main.MovieID = movieData[0].MovieID
					main.PathURL = "/movie/details/" + movieData[0].MovieID
				}

				if main.PathURL == "" {
					mu.Lock()
					firstErr = fiber.NewError(fiber.StatusNotFound, "Tidak ditemukan di MovieDetails")
					mu.Unlock()
					return
				}

				mu.Lock()
				results = append(results, main)
				pageRes = 1
				mu.Unlock()
			}()

		case "kdrama":
			wg.Add(1)
			go func() {
				defer wg.Done()
				tmp := *params

				movieData, _, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(&tmp))
				if err != nil {
					s.Log.Errorf("MovieSvc.GetAll error (kdrama): %v", err)
				}

				mdlResults, _, _, err := s.MdlSvc.GetAll(c, convert_types.MapToMdlQuery(&tmp))
				if err != nil || len(mdlResults) == 0 {
					s.Log.Errorf("MDL fetch error: %v", err)
					return
				}
				main := mdlResults[0]

				if len(movieData) > 0 {
					main.MovieID = movieData[0].MovieID
					main.PathURL = "/movie/details/" + movieData[0].MovieID
				}

				if main.PathURL == "" {
					mu.Lock()
					firstErr = fiber.NewError(fiber.StatusNotFound, "Tidak ditemukan di MovieDetails")
					mu.Unlock()
					return
				}

				mu.Lock()
				results = append(results, main)
				pageRes = 1
				mu.Unlock()
			}()
		}

		wg.Wait()
		if firstErr != nil && len(results) == 0 {
			return nil, 0, firstErr
		}

		return results, pageRes, nil
	}

	switch strings.ToLower(params.Type) {
	case "anime":
		wg.Add(1)
		go func() {
			defer wg.Done()

			data, page, err := s.AnilistSvc.GetAll(c, convert_types.MapToAnilistQuery(params))
			if err != nil {
				s.Log.Errorf("AnilistSvc.GetAll error: %v", err)
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			var filtered []response.MovieDetailOnlyResponse
			for _, d := range data {
				found := false

				tmp := *params
				tmp.Search = d.Title
				movieData, _, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(&tmp))
				if err == nil && len(movieData) > 0 {
					d.MovieID = movieData[0].MovieID
					d.PathURL = "/movie/details/" + movieData[0].MovieID
					filtered = append(filtered, d)
					found = true
				}

				if !found {
					odResults, err := s.OdService.GetAnimeByTitle(d.Title)
					if err == nil && len(odResults) > 0 {
						for _, od := range odResults {
							if od.URL != "" {
								d.MovieID = od.URL
								d.PathURL = "/otakudesu/detail/" + od.URL
								filtered = append(filtered, d)
								break
							}
						}
					}
				}
			}

			mu.Lock()
			results = append(results, filtered...)
			pageRes = page
			mu.Unlock()
		}()
	case "movie", "tv":
		wg.Add(1)
		go func() {
			defer wg.Done()

			data, page, err := s.TmdbSvc.GetAll(c, convert_types.MapToTmdbQuery(params))
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
			pageRes = page
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
					filtered = append(filtered, d)
				}
			}

			mu.Lock()
			results = append(results, filtered...)
			pageRes = page
			mu.Unlock()
		}()

	default:
		return nil, 0, fiber.NewError(fiber.StatusBadRequest, "Invalid query params")
	}

	wg.Wait()

	if firstErr != nil && len(results) == 0 {
		return nil, 0, firstErr
	}

	return results, pageRes, nil
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

	results, _, err := s.GetDiscover(c, params)
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
				detail.MovieID = odData[0].URL
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

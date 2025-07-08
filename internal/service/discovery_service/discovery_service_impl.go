package service

import (
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
		switch strings.ToLower(params.Type) {
		case "anime":
			wg.Add(1)
			go func() {
				defer wg.Done()

				tmp := *params
				movieData, page, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(&tmp))
				if err != nil {
					s.Log.Errorf("MovieSvc.GetAll error: %v", err)
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
					return
				}

				var filtered []response.MovieDetailOnlyResponse

				if len(movieData) > 0 {
					converted := convert_types.ConvertMvDetailToOnlyResp(movieData)
					for _, d := range converted {
						if d.MovieID != "" {
							filtered = append(filtered, d)
						}
					}
				} else {
					odResults, err := s.OdService.GetAnimeByTitle(tmp.Search)
					if err == nil && len(odResults) > 0 {
						for _, od := range odResults {
							if od.Title != "" && od.URL != "" {
								filtered = append(filtered, response.MovieDetailOnlyResponse{
									Title:        od.Title,
									MovieID:      od.URL,
									MovieType:    "anime",
									PathURL:      "/otakudesu/detail/" + od.URL,
									ThumbnailURL: od.ThumbnailURL,
									Genres: func() (r []string) {
										for _, g := range od.Genres {
											r = append(r, g.Title)
										}
										return
									}()})
							}
						}
					}
				}

				mu.Lock()
				results = append(results, filtered...)
				pageRes = page
				mu.Unlock()
			}()

		default:
			wg.Add(1)
			go func() {
				defer wg.Done()

				movieData, page, err := s.MovieSvc.GetAll(c, convert_types.MapToMovieDtQuery(params))
				if err != nil {
					s.Log.Errorf("MovieSvc.GetAll error: %v", err)
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
					return
				}

				movDataRes := convert_types.ConvertMvDetailToOnlyResp(movieData)

				mu.Lock()
				results = append(results, movDataRes...)
				pageRes = page
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
								d.PathURL = "/otakudesu/detail/ " + od.URL
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

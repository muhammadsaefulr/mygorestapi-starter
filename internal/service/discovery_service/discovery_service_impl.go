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
}

func NewDiscoveryService(validate *validator.Validate, SvcAn svcAnilist.AnilistServiceInterface, svcTmdb svcTmdb.TmdbServiceInterface, svcMdl svcMdl.MdlServiceInterface) DiscoveryServiceInterface {
	return &DiscoveryService{
		Log:        utils.Log,
		Validate:   validate,
		AnilistSvc: SvcAn,
		TmdbSvc:    svcTmdb,
		MdlSvc:     svcMdl,
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
		results  []response.MovieDetailOnlyResponse
		pageRes  int64
		firstErr error
		mu       sync.Mutex
		wg       sync.WaitGroup
	)

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
			mu.Lock()
			results = append(results, data...)
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
			mu.Lock()
			results = append(results, data...)
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
			mu.Lock()
			results = append(results, data...)
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

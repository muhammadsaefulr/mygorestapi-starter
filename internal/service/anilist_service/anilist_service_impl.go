package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/anilist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"

	// "github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/convert_types"
	modules "github.com/muhammadsaefulr/NimeStreamAPI/internal/infrastructure/modules/anilist"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/sirupsen/logrus"
)

type AnilistService struct {
	Log      *logrus.Logger
	Validate *validator.Validate
}

func NewAnilistService(validate *validator.Validate) AnilistServiceInterface {
	return &AnilistService{
		Log:      utils.Log,
		Validate: validate,
	}
}

func (s *AnilistService) GetAll(c *fiber.Ctx, params *request.QueryAnilist) ([]response.MovieDetailOnlyResponse, int64, error) {
	// Validasi input
	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	result, totalPage, err := modules.FetchAniListMedia(params.Category, params.Search, params)
	if err != nil {
		return nil, 0, err
	}

	return result, int64(totalPage), nil
}

func (s *AnilistService) GetMovieDetailsByID(c *fiber.Ctx, id string) (*response.MovieDetailOnlyResponse, error) {
	detail, err := modules.FetchAniListDetail(id)
	if err != nil {
		return nil, err
	}

	// log.Print("GetMovieDetailsByID id: ", id)
	// log.Printf("GetMovieDetailsByID detail: %+v", detail)

	return &detail, nil
}

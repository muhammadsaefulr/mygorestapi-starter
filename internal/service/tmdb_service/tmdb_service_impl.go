package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/tmdb/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/infrastructure/modules/tmdb"
)

type TMDbService struct {
	Validate *validator.Validate
}

func NewTMDbService(validate *validator.Validate) *TMDbService {
	return &TMDbService{Validate: validate}
}

func (s *TMDbService) GetAll(c *fiber.Ctx, params *request.QueryTmdb) ([]response.MovieDetailOnlyResponse, int64, error) {

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	mediaType := params.Type
	if mediaType != "tv" && mediaType != "movie" {
		mediaType = "movie"
	}

	result, err := modules.FetchTMDbMedia(params.Category, params.Search, mediaType, params)
	if err != nil {
		return nil, 0, err
	}

	var responseList []response.MovieDetailOnlyResponse
	for _, movie := range result {
		responseList = append(responseList, response.MovieDetailOnlyResponse{
			MovieID:      movie.MovieID,
			MovieType:    movie.MovieType,
			ThumbnailURL: movie.ThumbnailURL,
			Title:        movie.Title,
			Rating:       movie.Rating,
			Status:       movie.Status,
			TotalEps:     movie.TotalEps,
			Studio:       movie.Studio,
			ReleaseDate:  movie.ReleaseDate,
			Synopsis:     movie.Synopsis,
			Genres:       movie.Genres,
		})
	}

	return responseList, int64(len(responseList)), nil
}

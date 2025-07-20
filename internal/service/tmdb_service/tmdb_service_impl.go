package service

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/tmdb/request"
	modules "github.com/muhammadsaefulr/NimeStreamAPI/internal/infrastructure/modules/tmdb"
)

type TMDbService struct {
	Validate *validator.Validate
}

func NewTMDbService(validate *validator.Validate) *TMDbService {
	return &TMDbService{Validate: validate}
}

func (s *TMDbService) GetAll(c *fiber.Ctx, params *request.QueryTmdb) ([]response.MovieDetailOnlyResponse, int64, int64, error) {

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, 0, err
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

	if mediaType == "movie" && params.Category == "ongoing" {
		return nil, 0, 0, fiber.NewError(fiber.StatusBadRequest, "Ongoing Only available for tv show, drama or anime !")
	}

	result, totalPage, err := modules.FetchTMDbMedia(mediaType, params)
	if err != nil {
		return nil, 0, 0, err
	}

	var responseList []response.MovieDetailOnlyResponse
	for _, movie := range result {
		responseList = append(responseList, response.MovieDetailOnlyResponse{
			IDSource:     movie.IDSource,
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

	return responseList, int64(len(responseList)), int64(totalPage), nil
}

func (s *TMDbService) GetAllGenres(c *fiber.Ctx) ([]response.GenreDetail, error) {
	result, err := modules.GetTmdbListAllGenres()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *TMDbService) GetDetailByID(c *fiber.Ctx, id string, typeMov string) (*response.MovieDetailOnlyResponse, error) {
	idstr, errParseInt := strconv.Atoi(id)
	if errParseInt != nil {
		return &response.MovieDetailOnlyResponse{}, fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	results, err := modules.FetchTMDbDetail(idstr, typeMov, true)
	if err != nil {
		return &response.MovieDetailOnlyResponse{}, fiber.NewError(fiber.StatusInternalServerError, "Failed To Parsing ID")
	}

	return &results, nil
}

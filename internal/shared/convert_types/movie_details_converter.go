package convert_types

import (
	"fmt"
	"strings"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateMovieDetailsToModel(req *request.CreateMovieDetails) *model.MovieDetails {
	return &model.MovieDetails{
		MovieID:      req.MovieID,
		MovieType:    req.MovieType,
		ThumbnailURL: req.ThumbnailURL,
		Title:        req.Title,
		Rating:       req.Rating,
		Producer:     req.Producer,
		Status:       req.Status,
		TotalEps:     req.TotalEps,
		Studio:       req.Studio,
		ReleaseDate:  req.ReleaseDate,
		Synopsis:     req.Synopsis,
		Genres:       req.Genres,
	}
}

func UpdateMovieDetailsToModel(req *request.UpdateMovieDetails) *model.MovieDetails {
	return &model.MovieDetails{
		MovieID:      req.MovieID,
		MovieType:    req.MovieType,
		ThumbnailURL: req.ThumbnailURL,
		Title:        req.Title,
		Rating:       req.Rating,
		Producer:     req.Producer,
		Status:       req.Status,
		TotalEps:     req.TotalEps,
		Studio:       req.Studio,
		ReleaseDate:  req.ReleaseDate,
		Synopsis:     req.Synopsis,
		Genres:       req.Genres,
	}
}

func MovieDetailsModelToResp(
	movie *model.MovieDetails,
	rekomen *model.MovieDetails,
) *response.MovieDetailsResponse {

	episodesResp := make([]response.EpisodesResponse, 0, len(movie.Episodes))
	seen := make(map[string]bool)

	for _, ep := range movie.Episodes {
		if seen[ep.MovieEpsID] {
			continue
		}
		seen[ep.MovieEpsID] = true

		episodesResp = append(episodesResp, response.EpisodesResponse{
			MovieEpsId: ep.MovieEpsID,
			Title:      strings.Title(strings.ReplaceAll(strings.ReplaceAll(ep.MovieEpsID, "-eps-", "-episode-"), "-", " ")),
			VideoURL:   fmt.Sprintf("/movie/episodes/%s/%s", movie.MovieID, ep.MovieEpsID),
		})
	}

	convertMovie := func(m *model.MovieDetails) *response.MovieDetailOnlyResponse {
		return &response.MovieDetailOnlyResponse{
			MovieID:      m.MovieID,
			MovieType:    m.MovieType,
			ThumbnailURL: m.ThumbnailURL,
			Title:        m.Title,
			Rating:       m.Rating,
			Producer:     m.Producer,
			Status:       m.Status,
			TotalEps:     m.TotalEps,
			Studio:       m.Studio,
			ReleaseDate:  m.ReleaseDate,
			Synopsis:     m.Synopsis,
			Genres:       m.Genres,
			CreatedAt:    &m.CreatedAt,
			UpdatedAt:    &m.UpdatedAt,
		}
	}

	var rekomendResp *response.MovieDetailOnlyResponse
	if rekomen != nil {
		rekomendResp = convertMovie(rekomen)
	}

	return &response.MovieDetailsResponse{
		MovieDetail: convertMovie(movie),
		Episodes:    episodesResp,
		Rekomend:    rekomendResp,
	}
}

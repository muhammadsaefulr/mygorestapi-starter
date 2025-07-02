package convert_types

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/request"
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

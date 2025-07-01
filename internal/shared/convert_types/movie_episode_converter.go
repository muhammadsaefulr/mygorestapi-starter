package convert_types

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateMovieEpisodesToModel(req *request.CreateMovieEpisodes) *model.MovieEpisode {
	return &model.MovieEpisode{
		MovieEpsID: req.MovieEpsID,
		MovieId:    req.MovieId,
		Title:      req.Title,
		UploadedBy: req.UploadedBy,
		VideoURL:   req.ContentUploads,
	}
}

func UpdateMovieEpisodesToModel(req *request.UpdateMovieEpisodes) *model.MovieEpisode {
	return &model.MovieEpisode{
		MovieEpsID: req.MovieEpsID,
		MovieId:    req.MovieId,
		VideoURL:   req.ContentUploads,
	}
}

package convert_types

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/history/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/history/response"
	responses "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateHistoryToModel(req *request.CreateHistory) *model.History {
	return &model.History{
		UserId:       uuid.MustParse(req.UserId),
		MovieId:      req.MovieId,
		MovieEpsId:   req.MovieEpsId,
		PlaybackTime: req.PlaybackTime,
	}
}

func UpdateHistoryToModel(req *request.UpdateHistory) *model.History {
	return &model.History{
		ID:           req.ID,
		PlaybackTime: req.PlaybackTime,
	}
}

func HistoryToResponse(history *model.History, animeDetail *responses.MovieDetailOnlyResponse) response.HistoryResponse {
	return response.HistoryResponse{
		ID:           history.ID,
		UserId:       history.UserId.String(),
		MovieId:      history.MovieId,
		MovieEpsId:   history.MovieEpsId,
		PlaybackTime: history.PlaybackTime,
		AnimeDetail: responses.MovieDetailOnlyResponse{
			MovieID:      animeDetail.MovieID,
			MovieType:    animeDetail.MovieType,
			ThumbnailURL: animeDetail.ThumbnailURL,
			Title:        animeDetail.Title,
			Rating:       animeDetail.Rating,
			Producer:     animeDetail.Producer,
			Status:       animeDetail.Status,
			TotalEps:     animeDetail.TotalEps,
			Studio:       animeDetail.Studio,
			ReleaseDate:  animeDetail.ReleaseDate,
			Synopsis:     animeDetail.Synopsis,
			Genres:       animeDetail.Genres,
		},
		CreatedAt: history.CreatedAt,
		UpdatedAt: history.UpdatedAt,
	}
}

package convert_types

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/history/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateHistoryToModel(req *request.CreateHistory) *model.History {
	return &model.History{
		UserId:       uuid.MustParse(req.UserId),
		MovieEpsId:   req.MovieEpsId,
		PlaybackTime: req.PlaybackTime,
	}
}

func UpdateHistoryToModel(req *request.UpdateHistory) *model.History {
	return &model.History{
		ID:           req.ID,
		UserId:       uuid.MustParse(req.UserId),
		MovieEpsId:   req.MovieEpsId,
		PlaybackTime: req.PlaybackTime,
	}
}

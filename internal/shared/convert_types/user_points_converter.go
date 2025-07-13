package convert_types

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_points/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateUserPointsToModel(req *request.UserPoints) *model.UserPoints {
	return &model.UserPoints{
		UserID: uuid.MustParse(req.UserId),
		Value:  req.Value,
	}
}

func UpdateUserPointsToModel(req *request.UserPoints) *model.UserPoints {
	return &model.UserPoints{
		UserID: uuid.MustParse(req.UserId),
		Value:  req.Value,
	}
}

package convert_types

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/anilist/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateAnilistToModel(req *request.CreateAnilist) *model.Anilist {
	return &model.Anilist{
		// TODO: sesuaikan field sesuai model
		// Example:
		// Name: req.Name,
	}
}

func UpdateAnilistToModel(req *request.UpdateAnilist) *model.Anilist {
	return &model.Anilist{
		// TODO: sesuaikan field sesuai model
		// Example:
		// Name: req.Name,
	}
}

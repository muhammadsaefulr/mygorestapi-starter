package convert_types

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/mdl/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateMdlToModel(req *request.CreateMdl) *model.Mdl {
	return &model.Mdl{
		// TODO: sesuaikan field sesuai model
		// Example:
		// Name: req.Name,
	}
}

func UpdateMdlToModel(req *request.UpdateMdl) *model.Mdl {
	return &model.Mdl{
		// TODO: sesuaikan field sesuai model
		// Example:
		// Name: req.Name,
	}
}

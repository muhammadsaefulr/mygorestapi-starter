package convert_types

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/tmdb/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateTmdbToModel(req *request.CreateTmdb) *model.Tmdb {
	return &model.Tmdb{
		// TODO: sesuaikan field sesuai model
		// Example:
		// Name: req.Name,
	}
}

func UpdateTmdbToModel(req *request.UpdateTmdb) *model.Tmdb {
	return &model.Tmdb{
		// TODO: sesuaikan field sesuai model
		// Example:
		// Name: req.Name,
	}
}

package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/mdl/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type MdlRepo interface {
	GetAll(ctx context.Context, param *request.QueryMdl) ([]model.Mdl, int64, error)
	GetByID(ctx context.Context, id uint) (*model.Mdl, error)
	Create(ctx context.Context, data *model.Mdl) error
	Update(ctx context.Context, data *model.Mdl) error
	Delete(ctx context.Context, id uint) error
}

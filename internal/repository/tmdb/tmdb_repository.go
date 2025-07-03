package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/tmdb/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type TmdbRepo interface {
	GetAll(ctx context.Context, param *request.QueryTmdb) ([]model.Tmdb, int64, error)
	GetByID(ctx context.Context, id uint) (*model.Tmdb, error)
	Create(ctx context.Context, data *model.Tmdb) error
	Update(ctx context.Context, data *model.Tmdb) error
	Delete(ctx context.Context, id uint) error
}

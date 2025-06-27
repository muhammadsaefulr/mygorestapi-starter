package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_movie/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type RequestMovieRepo interface {
	GetAll(ctx context.Context, param *request.QueryRequestMovie) ([]model.RequestMovie, int64, error)
	GetByID(ctx context.Context, id uint) (*model.RequestMovie, error)
	Create(ctx context.Context, data *model.RequestMovie) error
	Update(ctx context.Context, data *model.RequestMovie) error
	Delete(ctx context.Context, id uint) error
}

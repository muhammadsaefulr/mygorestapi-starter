package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type MovieDetailsRepo interface {
	GetCountAll(ctx context.Context) (int64, error)
	GetAll(ctx context.Context, param *request.QueryMovieDetails) ([]model.MovieDetails, int64, error)
	GetByID(ctx context.Context, id string) (*model.MovieDetails, error)
	GetByIDPreEps(ctx context.Context, id string) (*model.MovieDetails, error)
	Create(ctx context.Context, data *model.MovieDetails) error
	Update(ctx context.Context, data *model.MovieDetails) (*model.MovieDetails, error)
	Delete(ctx context.Context, id string) error
}

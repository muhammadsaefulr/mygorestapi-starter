package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type MovieEpisodeRepo interface {
	GetAll(ctx context.Context, param *request.QueryMovieEpisode) ([]model.MovieEpisode, int64, error)
	GetByID(ctx context.Context, movie_eps_id string) (*model.MovieEpisode, error)
	GetByMovieID(ctx context.Context, movie_id string, patam *request.QueryMovieEpisode) ([]model.MovieEpisode, int64, error)
	Create(ctx context.Context, data *model.MovieEpisode) error
	Update(ctx context.Context, data *model.MovieEpisode) error
	Delete(ctx context.Context, movie_eps_id string) error
}

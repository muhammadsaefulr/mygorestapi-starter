package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type WatchlistRepo interface {
	GetAllWatchlist(ctx context.Context, param *request.QueryWatchlist, id_user string) ([]model.Watchlist, int64, error)
	GetWatchlistByID(ctx context.Context, movie_id string) (*model.Watchlist, error)
	CreateWatchlist(ctx context.Context, data *model.Watchlist) error
	UpdateWatchlist(ctx context.Context, data *model.Watchlist) error
	DeleteWatchlist(ctx context.Context, movie_id string) error
}

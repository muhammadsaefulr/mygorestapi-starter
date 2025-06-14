package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type WatchlistRepo interface {
	GetAllWatchlist(ctx context.Context, param *request.QueryWatchlist) ([]model.Watchlist, int64, error)
	GetWatchlistByID(ctx context.Context, id uint) (*model.Watchlist, error)
	CreateWatchlist(ctx context.Context, data *model.Watchlist) error
	UpdateWatchlist(ctx context.Context, data *model.Watchlist) error
	DeleteWatchlist(ctx context.Context, id uint) error
}

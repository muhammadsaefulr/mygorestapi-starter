package repository

import (
	"context"
	"log"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type WatchlistRepositoryImpl struct {
	DB *gorm.DB
}

func NewWatchlistRepositoryImpl(db *gorm.DB) WatchlistRepo {
	return &WatchlistRepositoryImpl{
		DB: db,
	}
}

func (r *WatchlistRepositoryImpl) GetAllWatchlist(ctx context.Context, param *request.QueryWatchlist, user_id string) ([]model.Watchlist, int64, error) {
	var data []model.Watchlist

	query := r.DB.WithContext(ctx).Model(&model.Watchlist{}).Where("user_id = ?", user_id)
	offset := (param.Page - 1) * param.Limit

	log.Printf("param: %+v", param)

	// log.Println(param)

	if param.Search != "" {
		searchLike := "%" + param.Search + "%"
		query = query.Where("name LIKE ?", searchLike)
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, int64(len(data)), nil
}

func (r *WatchlistRepositoryImpl) CreateWatchlist(ctx context.Context, data *model.Watchlist) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *WatchlistRepositoryImpl) UpdateWatchlist(ctx context.Context, data *model.Watchlist) error {
	return r.DB.WithContext(ctx).Where("movie_id = ? AND user_id = ?", data.MovieId, data.UserId).Updates(data).Error
}

func (r *WatchlistRepositoryImpl) DeleteWatchlist(ctx context.Context, movie_id string, user_id string) error {
	return r.DB.WithContext(ctx).Where("movie_id = ? AND user_id = ?", movie_id, user_id).Delete(&model.Watchlist{}).Error
}

package repository

import (
	"context"

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
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.Watchlist{}).Where("user_id = ?", user_id)
	offset := (param.Page - 1) * param.Limit

	// log.Println(param)

	if param.Search != "" {
		searchLike := "%" + param.Search + "%"
		query = query.Where("name LIKE ?", searchLike)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *WatchlistRepositoryImpl) GetWatchlistByID(ctx context.Context, id uint) (*model.Watchlist, error) {
	var data model.Watchlist
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *WatchlistRepositoryImpl) CreateWatchlist(ctx context.Context, data *model.Watchlist) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *WatchlistRepositoryImpl) UpdateWatchlist(ctx context.Context, data *model.Watchlist) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *WatchlistRepositoryImpl) DeleteWatchlist(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.Watchlist{}).Error
}

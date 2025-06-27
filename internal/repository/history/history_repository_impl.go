package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/history/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type HistoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewHistoryRepositoryImpl(db *gorm.DB) HistoryRepo {
	return &HistoryRepositoryImpl{
		DB: db,
	}
}

func (r *HistoryRepositoryImpl) GetAllByUserId(ctx context.Context, UserId string, param *request.QueryHistory) ([]model.History, int64, error) {
	var data []model.History
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.History{}).Where("user_id = ?", UserId)
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *HistoryRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.History, error) {
	var data model.History
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *HistoryRepositoryImpl) Create(ctx context.Context, data *model.History) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *HistoryRepositoryImpl) Update(ctx context.Context, data *model.History) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *HistoryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.History{}).Error
}

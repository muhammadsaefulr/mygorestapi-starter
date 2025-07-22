package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_vip/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type RequestVipRepositoryImpl struct {
	DB *gorm.DB
}

func NewRequestVipRepositoryImpl(db *gorm.DB) RequestVipRepo {
	return &RequestVipRepositoryImpl{
		DB: db,
	}
}

func (r *RequestVipRepositoryImpl) GetAll(ctx context.Context, param *request.QueryRequestVip) ([]model.RequestVip, int64, error) {
	var data []model.RequestVip
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.RequestVip{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *RequestVipRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.RequestVip, error) {
	var data model.RequestVip
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *RequestVipRepositoryImpl) Create(ctx context.Context, data *model.RequestVip) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *RequestVipRepositoryImpl) Update(ctx context.Context, data *model.RequestVip) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *RequestVipRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.RequestVip{}).Error
}

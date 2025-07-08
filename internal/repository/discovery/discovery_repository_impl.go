package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/discovery/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type DiscoveryRepositoryImpl struct {
	DB *gorm.DB
}

func NewDiscoveryRepositoryImpl(db *gorm.DB) DiscoveryRepo {
	return &DiscoveryRepositoryImpl{
		DB: db,
	}
}

func (r *DiscoveryRepositoryImpl) GetAll(ctx context.Context, param *request.QueryDiscovery) ([]model.Discovery, int64, error) {
	var data []model.Discovery
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.Discovery{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *DiscoveryRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.Discovery, error) {
	var data model.Discovery
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *DiscoveryRepositoryImpl) Create(ctx context.Context, data *model.Discovery) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *DiscoveryRepositoryImpl) Update(ctx context.Context, data *model.Discovery) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *DiscoveryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.Discovery{}).Error
}

package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/subscription_plan/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type SubscriptionPlanRepositoryImpl struct {
	DB *gorm.DB
}

func NewSubscriptionPlanRepositoryImpl(db *gorm.DB) SubscriptionPlanRepo {
	return &SubscriptionPlanRepositoryImpl{
		DB: db,
	}
}

func (r *SubscriptionPlanRepositoryImpl) GetAll(ctx context.Context, param *request.QuerySubscriptionPlan) ([]model.SubscriptionPlan, int64, error) {
	var data []model.SubscriptionPlan
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.SubscriptionPlan{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *SubscriptionPlanRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.SubscriptionPlan, error) {
	var data model.SubscriptionPlan
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *SubscriptionPlanRepositoryImpl) Create(ctx context.Context, data *model.SubscriptionPlan) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *SubscriptionPlanRepositoryImpl) Update(ctx context.Context, data *model.SubscriptionPlan) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *SubscriptionPlanRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.SubscriptionPlan{}).Error
}

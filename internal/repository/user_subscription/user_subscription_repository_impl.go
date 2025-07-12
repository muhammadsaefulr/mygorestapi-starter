package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_subscription/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type UserSubscriptionRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserSubscriptionRepositoryImpl(db *gorm.DB) UserSubscriptionRepo {
	return &UserSubscriptionRepositoryImpl{
		DB: db,
	}
}

func (r *UserSubscriptionRepositoryImpl) GetAll(ctx context.Context, param *request.QueryUserSubscription) ([]model.UserSubscription, int64, error) {
	var data []model.UserSubscription
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.UserSubscription{}).Preload("User").Preload("SubscriptionPlan")
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *UserSubscriptionRepositoryImpl) GetByUserID(ctx context.Context, id string) (*model.UserSubscription, error) {
	var data model.UserSubscription
	if err := r.DB.WithContext(ctx).Preload("User").Preload("SubscriptionPlan").First(&data, "user_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *UserSubscriptionRepositoryImpl) Create(ctx context.Context, data *model.UserSubscription) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *UserSubscriptionRepositoryImpl) UpdateByUserId(ctx context.Context, data *model.UserSubscription) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", data.UserID.String()).Updates(data).Error
}

func (r *UserSubscriptionRepositoryImpl) DeleteByUserId(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", id).Delete(&model.UserSubscription{}).Error
}

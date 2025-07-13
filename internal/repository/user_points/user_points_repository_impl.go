package repository

import (
	"context"

	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type UserPointsRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserPointsRepositoryImpl(db *gorm.DB) UserPointsRepo {
	return &UserPointsRepositoryImpl{
		DB: db,
	}
}

func (r *UserPointsRepositoryImpl) GetByUserID(ctx context.Context, id string) (*model.UserPoints, error) {
	var data model.UserPoints
	if err := r.DB.WithContext(ctx).First(&data, "user_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *UserPointsRepositoryImpl) Create(ctx context.Context, data *model.UserPoints) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *UserPointsRepositoryImpl) Update(ctx context.Context, data *model.UserPoints) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *UserPointsRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.UserPoints{}).Error
}

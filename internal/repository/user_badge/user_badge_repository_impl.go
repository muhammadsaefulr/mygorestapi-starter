package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_badge/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type UserBadgeRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserBadgeRepositoryImpl(db *gorm.DB) UserBadgeRepo {
	return &UserBadgeRepositoryImpl{
		DB: db,
	}
}

func (r *UserBadgeRepositoryImpl) GetAll(ctx context.Context, param *request.QueryUserBadge) ([]model.UserBadge, int64, error) {
	var data []model.UserBadge
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.UserBadge{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *UserBadgeRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.UserBadge, error) {
	var data model.UserBadge
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *UserBadgeRepositoryImpl) Create(ctx context.Context, data *model.UserBadge) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *UserBadgeRepositoryImpl) Update(ctx context.Context, data *model.UserBadge) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *UserBadgeRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.UserBadge{}).Error
}

func (r *UserBadgeRepositoryImpl) CreateUserBadgeInfo(ctx context.Context, data *model.UserBadgeInfo) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *UserBadgeRepositoryImpl) DeleteUserBadgeInfo(ctx context.Context, user_id string) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", user_id).Delete(&model.UserBadgeInfo{}).Error
}

func (r *UserBadgeRepositoryImpl) GetUserBadgeInfoByUserID(ctx context.Context, user_id string) ([]model.UserBadgeInfo, error) {
	var UserBadgeList []model.UserBadgeInfo

	err := r.DB.WithContext(ctx).
		Preload("Badge").
		Where("user_id = ?", user_id).
		Find(&UserBadgeList).Error

	return UserBadgeList, err
}

func (r *UserBadgeRepositoryImpl) UpdateUserBadgeInfo(ctx context.Context, data *model.UserBadgeInfo) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", data.UserID).Updates(data).Error
}

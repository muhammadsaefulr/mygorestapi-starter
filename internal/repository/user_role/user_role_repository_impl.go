package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_role/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type UserRoleRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRoleRepositoryImpl(db *gorm.DB) UserRoleRepo {
	return &UserRoleRepositoryImpl{
		DB: db,
	}
}

func (r *UserRoleRepositoryImpl) GetAll(ctx context.Context, param *request.QueryUserRole) ([]model.UserRole, int64, error) {
	var data []model.UserRole
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.UserRole{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *UserRoleRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.UserRole, error) {
	var data model.UserRole
	if err := r.DB.WithContext(ctx).Preload("Permissions").First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *UserRoleRepositoryImpl) Create(ctx context.Context, data *model.UserRole) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *UserRoleRepositoryImpl) Update(ctx context.Context, data *model.UserRole) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *UserRoleRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.UserRole{}).Error
}

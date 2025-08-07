package repository

import (
	"context"

	"github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/dto/role_permissions/request"
	model "github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/model"
	"gorm.io/gorm"
)

type RolePermissionsRepositoryImpl struct {
	DB *gorm.DB
}

func NewRolePermissionsRepositoryImpl(db *gorm.DB) RolePermissionsRepo {
	return &RolePermissionsRepositoryImpl{
		DB: db,
	}
}

func (r *RolePermissionsRepositoryImpl) GetAll(ctx context.Context, param *request.QueryRolePermissions) ([]model.RolePermissions, int64, error) {
	var data []model.RolePermissions
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.RolePermissions{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *RolePermissionsRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.RolePermissions, error) {
	var data model.RolePermissions
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *RolePermissionsRepositoryImpl) Create(ctx context.Context, data *model.RolePermissions) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *RolePermissionsRepositoryImpl) Update(ctx context.Context, data *model.RolePermissions) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *RolePermissionsRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.RolePermissions{}).Error
}

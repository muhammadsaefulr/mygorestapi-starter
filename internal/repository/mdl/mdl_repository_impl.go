package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/mdl/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type MdlRepositoryImpl struct {
	DB *gorm.DB
}

func NewMdlRepositoryImpl(db *gorm.DB) MdlRepo {
	return &MdlRepositoryImpl{
		DB: db,
	}
}

func (r *MdlRepositoryImpl) GetAll(ctx context.Context, param *request.QueryMdl) ([]model.Mdl, int64, error) {
	var data []model.Mdl
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.Mdl{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *MdlRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.Mdl, error) {
	var data model.Mdl
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *MdlRepositoryImpl) Create(ctx context.Context, data *model.Mdl) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *MdlRepositoryImpl) Update(ctx context.Context, data *model.Mdl) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *MdlRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.Mdl{}).Error
}

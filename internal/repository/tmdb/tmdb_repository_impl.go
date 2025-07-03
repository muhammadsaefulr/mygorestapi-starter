package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/tmdb/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type TmdbRepositoryImpl struct {
	DB *gorm.DB
}

func NewTmdbRepositoryImpl(db *gorm.DB) TmdbRepo {
	return &TmdbRepositoryImpl{
		DB: db,
	}
}

func (r *TmdbRepositoryImpl) GetAll(ctx context.Context, param *request.QueryTmdb) ([]model.Tmdb, int64, error) {
	var data []model.Tmdb
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.Tmdb{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *TmdbRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.Tmdb, error) {
	var data model.Tmdb
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *TmdbRepositoryImpl) Create(ctx context.Context, data *model.Tmdb) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *TmdbRepositoryImpl) Update(ctx context.Context, data *model.Tmdb) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *TmdbRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.Tmdb{}).Error
}

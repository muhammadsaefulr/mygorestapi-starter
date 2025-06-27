package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_movie/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type RequestMovieRepositoryImpl struct {
	DB *gorm.DB
}

func NewRequestMovieRepositoryImpl(db *gorm.DB) RequestMovieRepo {
	return &RequestMovieRepositoryImpl{
		DB: db,
	}
}

func (r *RequestMovieRepositoryImpl) GetAll(ctx context.Context, param *request.QueryRequestMovie) ([]model.RequestMovie, int64, error) {
	var data []model.RequestMovie
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.RequestMovie{}).Preload("RequestedBy")
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	if param.Search != "" {
		searchPattern := "%" + param.Search + "%"
		query = query.Where("title ILIKE ? OR genre ILIKE ? OR type_movie ILIKE ? OR CAST(user_id_request AS TEXT) ILIKE ?", searchPattern, searchPattern, searchPattern, searchPattern)
	}

	return data, total, nil
}

func (r *RequestMovieRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.RequestMovie, error) {
	var data model.RequestMovie
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Preload("RequestedBy").Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *RequestMovieRepositoryImpl) Create(ctx context.Context, data *model.RequestMovie) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *RequestMovieRepositoryImpl) Update(ctx context.Context, data *model.RequestMovie) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *RequestMovieRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.RequestMovie{}).Error
}

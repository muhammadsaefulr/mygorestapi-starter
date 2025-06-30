package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type MovieDetailsRepositoryImpl struct {
	DB *gorm.DB
}

func NewMovieDetailsRepositoryImpl(db *gorm.DB) MovieDetailsRepo {
	return &MovieDetailsRepositoryImpl{
		DB: db,
	}
}

func (r *MovieDetailsRepositoryImpl) GetAll(ctx context.Context, param *request.QueryMovieDetails) ([]model.MovieDetails, int64, error) {
	var data []model.MovieDetails
	var total int64
	query := r.DB.WithContext(ctx).Model(&model.MovieDetails{})

	if param.Search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?) ", "%"+param.Search+"%")
	}

	if param.Type != "" {
		query = query.Where("LOWER(movie_type) LIKE LOWER(?)", "%"+param.Type+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (param.Page - 1) * param.Limit
	if err := query.
		Limit(param.Limit).
		Offset(offset).
		Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *MovieDetailsRepositoryImpl) GetByID(ctx context.Context, id string) (*model.MovieDetails, error) {
	var data model.MovieDetails
	if err := r.DB.WithContext(ctx).First(&data, "movie_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *MovieDetailsRepositoryImpl) Create(ctx context.Context, data *model.MovieDetails) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *MovieDetailsRepositoryImpl) Update(ctx context.Context, data *model.MovieDetails) error {
	return r.DB.WithContext(ctx).Where("movie_id = ?", data.MovieID).Updates(data).Error
}

func (r *MovieDetailsRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Where("movie_id = ?", id).Delete(&model.MovieDetails{}).Error
}

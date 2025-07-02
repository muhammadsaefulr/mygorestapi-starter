package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type MovieEpisodeRepositoryImpl struct {
	DB *gorm.DB
}

func NewMovieEpisodeRepositoryImpl(db *gorm.DB) MovieEpisodeRepo {
	return &MovieEpisodeRepositoryImpl{
		DB: db,
	}
}

func (r *MovieEpisodeRepositoryImpl) GetAll(ctx context.Context, param *request.QueryMovieEpisode) ([]model.MovieEpisode, int64, error) {
	var data []model.MovieEpisode
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.MovieEpisode{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *MovieEpisodeRepositoryImpl) GetByID(ctx context.Context, movie_eps_id string) (*model.MovieEpisode, error) {
	var data model.MovieEpisode
	if err := r.DB.WithContext(ctx).First(&data, "movie_eps_id = ?", movie_eps_id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *MovieEpisodeRepositoryImpl) GetByMovieID(ctx context.Context, movie_id string) ([]model.MovieEpisode, error) {
	var data []model.MovieEpisode
	if err := r.DB.WithContext(ctx).Where("movie_id = ?", movie_id).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *MovieEpisodeRepositoryImpl) Create(ctx context.Context, data *model.MovieEpisode) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *MovieEpisodeRepositoryImpl) Update(ctx context.Context, data *model.MovieEpisode) error {

	return r.DB.WithContext(ctx).Model(&model.MovieEpisode{}).Where("movie_eps_id = ?", data.MovieEpsID).Updates(data).Error

}

func (r *MovieEpisodeRepositoryImpl) Delete(ctx context.Context, movie_eps_id string) error {
	return r.DB.WithContext(ctx).Where("movie_eps_id = ?", movie_eps_id).Delete(&model.MovieEpisode{}).Error
}

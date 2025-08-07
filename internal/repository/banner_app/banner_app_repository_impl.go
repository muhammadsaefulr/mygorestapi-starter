package repository

import (
	"context"
	"math"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/banner_app/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type BannerAppRepositoryImpl struct {
	DB *gorm.DB
}

func NewBannerAppRepositoryImpl(db *gorm.DB) BannerAppRepo {
	return &BannerAppRepositoryImpl{
		DB: db,
	}
}

func (r *BannerAppRepositoryImpl) GetAll(ctx context.Context, param *request.QueryBannerApp) ([]model.BannerApp, int64, error) {
	var (
		data  []model.BannerApp
		total int64
	)

	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}
	offset := (param.Page - 1) * param.Limit

	baseQuery := r.DB.WithContext(ctx).Model(&model.BannerApp{})

	if param.Type != "" {
		baseQuery = baseQuery.Where("banner_type = ?", param.Type)
	}

	countQuery := baseQuery.Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := baseQuery.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	totalPages := int64(math.Ceil(float64(total) / float64(param.Limit)))

	return data, totalPages, nil
}

func (r *BannerAppRepositoryImpl) GetByID(ctx context.Context, id uint) (*model.BannerApp, error) {
	var data model.BannerApp
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *BannerAppRepositoryImpl) Create(ctx context.Context, data *model.BannerApp) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *BannerAppRepositoryImpl) Update(ctx context.Context, data *model.BannerApp) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *BannerAppRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.BannerApp{}).Error
}

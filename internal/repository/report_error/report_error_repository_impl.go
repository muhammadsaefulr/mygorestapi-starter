package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/report_error/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type ReportErrorRepositoryImpl struct {
	DB *gorm.DB
}

func NewReportErrorRepositoryImpl(db *gorm.DB) ReportErrorRepo {
	return &ReportErrorRepositoryImpl{
		DB: db,
	}
}

func (r *ReportErrorRepositoryImpl) GetAll(ctx context.Context, param *request.QueryReportError) ([]model.ReportError, int64, error) {
	var data []model.ReportError
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.ReportError{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *ReportErrorRepositoryImpl) GetByID(ctx context.Context, id string) (*model.ReportError, error) {
	var data model.ReportError
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *ReportErrorRepositoryImpl) Create(ctx context.Context, data *model.ReportError) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *ReportErrorRepositoryImpl) Update(ctx context.Context, data *model.ReportError) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ReportId).Updates(data).Error
}

func (r *ReportErrorRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.ReportError{}).Error
}

package repository

import (
	"context"

	"{{.ModulePath}}/internal/domain/dto/{{.Name}}/request"
	model "{{.ModulePath}}/internal/domain/model"
	"gorm.io/gorm"
)

type {{.PascalName}}RepositoryImpl struct {
	DB *gorm.DB
}

func New{{.PascalName}}RepositoryImpl(db *gorm.DB) {{.PascalName}}Repo {
	return &{{.PascalName}}RepositoryImpl{
		DB: db,
	}
}

func (r *{{.PascalName}}RepositoryImpl) GetAll{{.PascalName}}(ctx context.Context, param *request.Query{{.PascalName}}) ([]model.{{.PascalName}}, int64, error) {
	var data []model.{{.PascalName}}
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.{{.PascalName}}{})
	offset := (param.Page - 1) * param.Limit

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(param.Limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (r *{{.PascalName}}RepositoryImpl) Get{{.PascalName}}ByID(ctx context.Context, id string) (*model.{{.PascalName}}, error) {
	var data model.{{.PascalName}}
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *{{.PascalName}}RepositoryImpl) Create{{.PascalName}}(ctx context.Context, data *model.{{.PascalName}}) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *{{.PascalName}}RepositoryImpl) Update{{.PascalName}}(ctx context.Context, data *model.{{.PascalName}}) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *{{.PascalName}}RepositoryImpl) Delete{{.PascalName}}(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.{{.PascalName}}{}).Error
}
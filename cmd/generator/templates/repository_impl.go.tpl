package repository

import (
	"context"
	"strings"

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

func (r *{{.PascalName}}RepositoryImpl) GetAll(ctx context.Context, param *request.Query{{.PascalName}}) ([]model.{{.PascalName}}, int64, int64, error) {
	var data []model.{{.PascalName}}
	var total int64
	var totalPages int64

	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}

	query := r.DB.WithContext(ctx).Model(&model.{{.PascalName}}{})

	if param.Search != "" {
		var stringColumns []string
		{{- range .Fields }}
		{{- if eq .Type "string" }}
		stringColumns = append(stringColumns, "COALESCE({{.SnakeCaseName}}, '')")
		{{- end }}
		{{- end }}
		
		if len(stringColumns) > 0 {
			document := strings.Join(stringColumns, " || ' ' || ")
			searchQuery := "to_tsvector('english', " + document + ") @@ plainto_tsquery('english', ?)"
			query = query.Where(searchQuery, param.Search)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	if param.Limit > 0 {
		totalPages = (total + int64(param.Limit) - 1) / int64(param.Limit)
	}

	offset := (param.Page - 1) * param.Limit
	err := query.Limit(param.Limit).Offset(offset).Order("id DESC").Find(&data).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return data, total, totalPages, nil
}

func (r *{{.PascalName}}RepositoryImpl) GetByID(ctx context.Context, id uint) (*model.{{.PascalName}}, error) {
	var data model.{{.PascalName}}
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *{{.PascalName}}RepositoryImpl) Create(ctx context.Context, data *model.{{.PascalName}}) error {
	return r.DB.WithContext(ctx).Create(data).Error
}

func (r *{{.PascalName}}RepositoryImpl) Update(ctx context.Context, data *model.{{.PascalName}}) error {
	return r.DB.WithContext(ctx).Where("id = ?", data.ID).Updates(data).Error
}

func (r *{{.PascalName}}RepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.{{.PascalName}}{}).Error
}

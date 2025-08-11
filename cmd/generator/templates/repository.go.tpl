package repository

import (
	"context"

	"{{.ModulePath}}/internal/domain/dto/{{.Name}}/request"
	model "{{.ModulePath}}/internal/domain/model"
)

type {{.PascalName}}Repo interface {
	GetAll(ctx context.Context, param *request.Query{{.PascalName}}) ([]model.{{.PascalName}}, int64, int64, error)
	GetByID(ctx context.Context, id uint) (*model.{{.PascalName}}, error)
	Create(ctx context.Context, data *model.{{.PascalName}}) error
	Update(ctx context.Context, data *model.{{.PascalName}}) error
	Delete(ctx context.Context, id uint) error
}

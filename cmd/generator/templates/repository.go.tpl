package repository

import (
	"context"

	"{{ .ModulePath }}/internal/domain/model/{{.Name}}"
	"{{ .ModulePath }}/internal/domain/dto/{{.Name}}/request"
)

type {{.PascalName}}Repo interface {
	GetAll{{.PascalName}}(ctx context.Context, param *request.Query{{.PascalName}}) ([]model.{{.PascalName}}, int64, error)
	Get{{.PascalName}}ByID(ctx context.Context, id string) (*model.{{.PascalName}}, error)
	Create{{.PascalName}}(ctx context.Context, data *model.{{.PascalName}}) error
	Update{{.PascalName}}(ctx context.Context, data *model.{{.PascalName}}) error
	Delete{{.PascalName}}(ctx context.Context, id string) error
}

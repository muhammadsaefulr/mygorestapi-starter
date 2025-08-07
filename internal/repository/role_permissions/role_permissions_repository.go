package repository

import (
	"context"

	"github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/dto/role_permissions/request"
	model "github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/model"
)

type RolePermissionsRepo interface {
	GetAll(ctx context.Context, param *request.QueryRolePermissions) ([]model.RolePermissions, int64, error)
	GetByID(ctx context.Context, id uint) (*model.RolePermissions, error)
	Create(ctx context.Context, data *model.RolePermissions) error
	Update(ctx context.Context, data *model.RolePermissions) error
	Delete(ctx context.Context, id uint) error
}

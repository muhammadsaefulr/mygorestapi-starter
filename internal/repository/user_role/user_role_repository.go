package repository

import (
	"context"

	"github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/dto/user_role/request"
	model "github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/model"
)

type UserRoleRepo interface {
	GetAll(ctx context.Context, param *request.QueryUserRole) ([]model.UserRole, int64, error)
	GetByID(ctx context.Context, id uint) (*model.UserRole, error)
	Create(ctx context.Context, data *model.UserRole) error
	Update(ctx context.Context, data *model.UserRole) error
	Delete(ctx context.Context, id uint) error
}

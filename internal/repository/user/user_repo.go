package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model/user"
)

type UserRepo interface {
	GetAllUser(ctx context.Context, param *request.QueryUser) ([]model.User, int64, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
}

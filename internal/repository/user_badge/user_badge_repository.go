package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_badge/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type UserBadgeRepo interface {
	GetAll(ctx context.Context, param *request.QueryUserBadge) ([]model.UserBadge, int64, error)
	GetByID(ctx context.Context, id uint) (*model.UserBadge, error)
	Create(ctx context.Context, data *model.UserBadge) error
	Update(ctx context.Context, data *model.UserBadge) error
	Delete(ctx context.Context, id uint) error

	GetUserBadgeInfoByUserID(ctx context.Context, user_id string) ([]model.UserBadgeInfo, error)
	CreateUserBadgeInfo(ctx context.Context, data *model.UserBadgeInfo) error
	UpdateUserBadgeInfo(ctx context.Context, data *model.UserBadgeInfo) error
	DeleteUserBadgeInfo(ctx context.Context, user_id string) error
}

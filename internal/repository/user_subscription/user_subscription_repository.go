package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_subscription/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type UserSubscriptionRepo interface {
	GetAll(ctx context.Context, param *request.QueryUserSubscription) ([]model.UserSubscription, int64, error)
	GetByUserID(ctx context.Context, id string) (*model.UserSubscription, error)
	Create(ctx context.Context, data *model.UserSubscription) error
	UpdateByUserId(ctx context.Context, data *model.UserSubscription) error
	DeleteByUserId(ctx context.Context, id string) error
}

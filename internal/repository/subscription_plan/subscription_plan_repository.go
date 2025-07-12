package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/subscription_plan/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type SubscriptionPlanRepo interface {
	GetAll(ctx context.Context, param *request.QuerySubscriptionPlan) ([]model.SubscriptionPlan, int64, error)
	GetByID(ctx context.Context, id uint) (*model.SubscriptionPlan, error)
	Create(ctx context.Context, data *model.SubscriptionPlan) error
	Update(ctx context.Context, data *model.SubscriptionPlan) error
	Delete(ctx context.Context, id uint) error
}

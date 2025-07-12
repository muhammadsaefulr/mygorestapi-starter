package convert_types

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user_subscription/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateUserSubscriptionToModel(req *request.CreateUserSubscription) *model.UserSubscription {
	return &model.UserSubscription{
		UserID:             uuid.MustParse(req.UserId),
		SubscriptionPlanID: req.SubscriptionPlanId,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		IsActive:           req.IsActive,
	}
}

func UpdateUserSubscriptionToModel(req *request.UpdateUserSubscription) *model.UserSubscription {
	return &model.UserSubscription{
		ID:                 req.ID,
		SubscriptionPlanID: req.SubscriptionPlanId,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		IsActive:           req.IsActive,
	}
}

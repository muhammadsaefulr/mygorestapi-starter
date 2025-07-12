package convert_types

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/subscription_plan/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateSubscriptionPlanToModel(req *request.CreateSubscriptionPlan) *model.SubscriptionPlan {
	return &model.SubscriptionPlan{
		PlanName: req.Name,
		Duration: req.Duration,
		Benefits: req.Benefits,
		Price:    req.Price,
	}
}

func UpdateSubscriptionPlanToModel(req *request.UpdateSubscriptionPlan) *model.SubscriptionPlan {
	return &model.SubscriptionPlan{
		PlanName: req.Name,
		Duration: req.Duration,
		Benefits: req.Benefits,
		Price:    req.Price,
	}
}

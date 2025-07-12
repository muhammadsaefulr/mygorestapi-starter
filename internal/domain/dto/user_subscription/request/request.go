package request

import "time"

type CreateUserSubscription struct {
	UserId             string    `json:"user_id" example:"7e26e5c0-0cb7-4d3f-a777-2a5adbb099fa"`
	SubscriptionPlanId uint      `json:"subscription_plan_id" example:"1"`
	StartDate          time.Time `json:"start_date" example:"2025-07-10T00:00:00Z"`
	EndDate            time.Time `json:"end_date" example:"2025-08-10T00:00:00Z"`
	IsActive           bool      `json:"is_active" example:"true"`
}

type UpdateUserSubscription struct {
	ID                 uint      `json:"-"`
	UserId             string    `json:"user_id" example:"7e26e5c0-0cb7-4d3f-a777-2a5adbb099fa"`
	SubscriptionPlanId uint      `json:"subscription_plan_id" example:"1"`
	StartDate          time.Time `json:"start_date" example:"2025-07-10T00:00:00Z"`
	EndDate            time.Time `json:"end_date" example:"2025-08-10T00:00:00Z"`
	IsActive           bool      `json:"is_active" example:"true"`
}

type QueryUserSubscription struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

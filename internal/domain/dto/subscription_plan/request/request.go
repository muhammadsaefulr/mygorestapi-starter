package request

type CreateSubscriptionPlan struct {
	Name     string   `json:"plan_name"`
	Duration int      `json:"duration"`
	Benefits []string `json:"benefits"`
	Price    int      `json:"price"`
}

type UpdateSubscriptionPlan struct {
	Name     string   `json:"plan_name"`
	Duration int      `json:"duration"`
	Benefits []string `json:"benefits"`
	Price    int      `json:"price"`
}

type QuerySubscriptionPlan struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

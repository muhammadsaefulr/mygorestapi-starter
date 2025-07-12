package response

import "time"

type SubscriptionPlanResponse struct {
    // contoh field, sesuaikan dengan kebutuhan
	ID            uint      `json:"id"`
	Name		  string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

package response

import "time"

type RequestVipResponse struct {
	ID            uint      `json:"id"`
	UserId        string    `json:"user_id"`
	Name          string    `json:"atas_nama_tf"`
	Email         string    `json:"email"`
	BuktiTf       string    `json:"bukti_tf"`
	PaymentMethod string    `json:"payment_method"`
	StatusAcc     string    `json:"status_acc"`
	UpdatedBy     string    `json:"updated_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

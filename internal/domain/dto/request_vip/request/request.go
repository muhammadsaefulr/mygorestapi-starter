package request

import "mime/multipart"

type CreateRequestVip struct {
	UserId        string                `json:"-"`
	PaymentMethod string                `json:"payment_method"`
	Name          string                `json:"atas_nama_tf"`
	Email         string                `json:"email"`
	BuktiTf       *multipart.FileHeader `json:"-"`
	BuktiTfStr    string                `json:"-"`
	StatusAcc     string                `json:"-"`
}

type UpdateRequestVip struct {
	ID        uint   `json:"-"`
	Name      string `json:"atas_nama_tf"`
	Email     string `json:"email"`
	BuktiTf   string `json:"bukti_tf"`
	StatusAcc string `json:"status_acc"`
	UpdatedBy string `json:"-"`
}

type QueryRequestVip struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

package response

import "time"

type {{.PascalName}}Response struct {
    // contoh field, sesuaikan dengan kebutuhan
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

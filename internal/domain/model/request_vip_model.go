package model

import (
	"time"

	"github.com/google/uuid"
)

type RequestVip struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"not null" json:"user_id"`
	Email         string    `gorm:"not null" json:"email"`
	PaymentMethod string    `gorm:"not null" json:"payment_method"`
	AtasNamaTf    string    `gorm:"not null" json:"atas_nama_tf"`
	BuktiTf       string    `gorm:"not null" json:"bukti_tf"`
	StatusAcc     string    `gorm:"not null" json:"status_acc"`
	UpdatedBy     uuid.UUID `gorm:"not null" json:"updated_by"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type RequestVip struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uuid.UUID `gorm:"not null"`
	Email         string    `gorm:"not null"`
	PaymentMethod string    `gorm:"not null"`
	AtasNamaTf    string    `gorm:"not null"`
	BuktiTf       string    `gorm:"not null"`
	StatusAcc     string    `gorm:"not null"`
	UpdatedBy     uuid.UUID `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

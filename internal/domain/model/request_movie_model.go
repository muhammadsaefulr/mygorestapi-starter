package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RequestMovie struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserIdRequest uuid.UUID      `gorm:"not null" json:"user_id_request"`
	Title         string         `gorm:"not null" json:"title"`
	TypeMovie     string         `gorm:"not null" json:"type_movie"`
	Genre         pq.StringArray `gorm:"type:text[];not null" json:"genre"`
	Description   string         `gorm:"not null" json:"description"`
	StatusMovie   string         `gorm:"not null" json:"status_movie"`
	StatusRequest string         `gorm:"not null" json:"status_request"`
	RequestedBy   User           `gorm:"foreignKey:UserIdRequest;references:ID" json:"requested_by"`
	CreatedAt     time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

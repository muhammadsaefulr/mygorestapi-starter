package model

import (
	"time"

	"github.com/google/uuid"
)

type UserBadge struct {
	ID        uint   `gorm:"primaryKey"`
	BadgeName string `gorm:"not null;unique"`
	IconURL   string
	Color     string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type UserBadgeInfo struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	BadgeID   uint      `gorm:"not null;uniqueIndex"`
	GivenAt   time.Time `gorm:"autoCreateTime"`
	Note      string    `json:"note"`
	HandledBy uuid.UUID `gorm:"type:uuid;not null"`

	Badge UserBadge `gorm:"foreignKey:BadgeID;references:ID"`
}

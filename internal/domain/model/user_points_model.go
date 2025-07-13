package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPoints struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	UserID    uuid.UUID `gorm:"not null;uniqueIndex" json:"user_id"`
	Value     int       `gorm:"not null" json:"value"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (u *UserPoints) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

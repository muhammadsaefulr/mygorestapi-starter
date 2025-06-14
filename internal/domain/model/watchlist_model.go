package model

import (
	"time"

	"github.com/google/uuid"
)

type Watchlist struct {
	ID        uint      `gorm:"primaryKey"`
	UserId    uuid.UUID `gorm:"not null" json:"user_id"`
	MovieId   string    `gorm:"not null" json:"movie_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

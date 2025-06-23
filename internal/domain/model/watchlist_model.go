package model

import (
	"time"

	"github.com/google/uuid"
)

type Watchlist struct {
	ID            uint      `gorm:"primaryKey"`
	UserId        uuid.UUID `gorm:"not null;index:idx_user_movie,unique"`
	MovieId       string    `gorm:"not null;index:idx_user_movie,unique"`
	ThumbImageUrl string    `gorm:"not null" json:"thumb_image_url"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

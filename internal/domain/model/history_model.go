package model

import (
	"time"

	"github.com/google/uuid"
)

type History struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	UserId       uuid.UUID   `gorm:"not null" json:"user_id"`
	MovieId      string      `gorm:"not null" json:"movie_id"`
	MovieEpsId   string      `gorm:"uniqueIndex:idx_history_movie_eps_id" json:"movie_eps_id"`
	PlaybackTime int         `gorm:"not null" json:"playback_time"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	DetailMovie  AnimeDetail `gorm:"-"`
	UpdatedAt    time.Time
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type TrackEpisodeView struct {
	ID             uint      `gorm:"primaryKey"`
	UserId         uuid.UUID `gorm:"not null" json:"user_id"`
	MovieDetailUrl string    `gorm:"not null" json:"movie_detail_url"`
	EpisodeId      string    `gorm:"not null" json:"episode_id"`
	WatchedAt      time.Time `gorm:"autoCreateTime:milli" json:"watched_at"`
	CreatedAt      time.Time `gorm:"autoCreateTime:milli" json:"-"`
}

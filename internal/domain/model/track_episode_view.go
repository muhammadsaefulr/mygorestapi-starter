package model

import (
	"time"

	"github.com/google/uuid"
)

type TrackEpisodeView struct {
	ID             uint      `gorm:"primaryKey"`
	UserId         uuid.UUID `gorm:"not null" json:"user_id"`
	TypeMovie      string    `gorm:"column:type_movie;not null;default:'anime'" json:"type_movie"`
	MovieDetailUrl string    `gorm:"column:movie_detail_url" json:"movie_detail_url"`
	EpisodeId      string    `gorm:"not null" json:"episode_id"`
	WatchedAt      time.Time `gorm:"autoCreateTime:milli" json:"watched_at"`
	CreatedAt      time.Time `gorm:"autoCreateTime:milli" json:"-"`
}

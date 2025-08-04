package model

import (
	"time"
)

type MovieEpisode struct {
	ID         uint   `gorm:"primaryKey"`
	MovieEpsID string `gorm:"not null;index:idx_eps_unique,unique" json:"movie_eps_id"`
	MovieId    string `gorm:"not null;index:idx_eps_unique,unique" json:"movie_id"`
	Resolution string `gorm:"not null;index:idx_eps_unique,unique" json:"resolution"`
	VideoURL   string `gorm:"not null;index:idx_eps_unique,unique" json:"video_url"`
	Title      string `gorm:"not null" json:"title"` // title untuk penentu source nya misal pixeldrain atau uploads

	SourceBy  string    `gorm:"not null" json:"source_by"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at,omitempty"`
}

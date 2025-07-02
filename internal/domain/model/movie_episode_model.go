package model

import (
	"time"
)

type MovieEpisode struct {
	ID         uint   `gorm:"primaryKey"`
	MovieEpsID string `gorm:"not null;index:idx_eps_unique,unique"`
	MovieId    string `gorm:"not null;index:idx_eps_unique,unique"`
	Resolution string `gorm:"not null;index:idx_eps_unique,unique"`
	VideoURL   string `gorm:"not null;index:idx_eps_unique,unique"`
	Title      string `gorm:"not null"` // title untuk penentu source nya misal pixeldrain atau uploads

	SourceBy  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

package model

import (
	"time"
)

type BannerApp struct {
	ID         uint      `gorm:"primaryKey"`
	Title      string    `gorm:"not null;uniqueIndex"`
	ImageUrl   string    `gorm:"not null"`
	BannerType string    `gorm:"not null;uniqueIndex"` // e.g., "movie", "drakor", "anime"
	DetailURL  string    `gorm:"not null;uniqueIndex"`
	UpdatedBy  string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

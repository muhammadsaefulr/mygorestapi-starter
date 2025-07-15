package model

import (
	"time"
)

type BannerApp struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	ImageUrl  string    `gorm:"not null"`
	DetailURL string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

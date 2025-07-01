package model

import (
	"time"
)

type MovieEpisode struct {
	MovieEpsID string    `gorm:"primaryKey"`
	MovieId    string    `gorm:"unique;not null"`
	Title      string    `gorm:"not null"`
	VideoURL   string    `gorm:"not null"`
	Resolution string    `gorm:"not null"`
	UploadedBy string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

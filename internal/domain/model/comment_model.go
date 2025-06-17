package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey;not_null" json:"id"`
	UserId    uuid.UUID `gorm:"not null" json:"user_id"`
	MovieId   string    `gorm:"not null" json:"movie_id"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time
}

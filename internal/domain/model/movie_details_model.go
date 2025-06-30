package model

import (
	"time"

	"github.com/lib/pq"
)

type MovieDetails struct {
	MovieID      string         `gorm:"primaryKey" json:"movie_id"`
	MovieType    string         `gorm:"type:varchar(50)" json:"movie_type"`
	ThumbnailURL string         `gorm:"type:text" json:"thumbnail_url"`
	Title        string         `gorm:"type:varchar(255)" json:"title"`
	Rating       string         `gorm:"type:varchar(50)" json:"rating"`
	Producer     string         `gorm:"type:varchar(255)" json:"producer"`
	Status       string         `gorm:"type:varchar(50)" json:"status"`
	TotalEps     string         `gorm:"type:varchar(50)" json:"total_eps"`
	Studio       string         `gorm:"type:varchar(255)" json:"studio"`
	ReleaseDate  string         `gorm:"type:varchar(100)" json:"release_date"`
	Synopsis     string         `gorm:"type:text" json:"synopsis"`
	Genres       pq.StringArray `gorm:"type:text[]" json:"genres"`

	Episodes []MovieUploader `gorm:"foreignKey:MovieId;references:MovieID" json:"episodes,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

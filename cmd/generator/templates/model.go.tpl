package model

import (
	"time"

)

type {{.PascalName}} struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

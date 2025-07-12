package model

import (
	"time"

	"github.com/lib/pq"
)

type SubscriptionPlan struct {
	ID       uint           `gorm:"primaryKey"`
	PlanName string         `gorm:"not null;unique"`
	Duration int            `gorm:"not null"`
	Price    int            `gorm:"not null"`
	Benefits pq.StringArray `gorm:"type:text[]" json:"benefit"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

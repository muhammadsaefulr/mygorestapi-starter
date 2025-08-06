package model

import (
	"time"

	"github.com/lib/pq"
)

type SubscriptionPlan struct {
	ID       uint           `gorm:"primaryKey" json:"id"`
	PlanName string         `gorm:"not null;unique" json:"plan_name"`
	Duration int            `gorm:"not null" json:"duration"`
	Price    int            `gorm:"not null" json:"price"`
	Benefits pq.StringArray `gorm:"type:text[]" json:"benefit"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

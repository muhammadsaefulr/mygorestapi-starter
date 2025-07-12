package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSubscription struct {
	ID uint `gorm:"primaryKey"`

	UserID             uuid.UUID `gorm:"not null;unique"`
	SubscriptionPlanID uint      `gorm:"not null"`

	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	IsActive  bool      `gorm:"default:true"`
	UpdatedBy uuid.UUID `gorm:"not null"`

	User             *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SubscriptionPlan *SubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID" json:"subscription_plan,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

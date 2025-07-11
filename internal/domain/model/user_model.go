package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID         `gorm:"primaryKey;not null" json:"id"`
	Name             string            `gorm:"not null" json:"name"`
	Email            string            `gorm:"uniqueIndex;not null" json:"email"`
	Password         string            `gorm:"not null" json:"-"`
	Role             string            `gorm:"default:user;not null" json:"role"`
	RoleId           uint              `gorm:"default:1;not null" json:"role_id"`
	VerifiedEmail    bool              `gorm:"default:false;not null" json:"verified_email"`
	CreatedAt        time.Time         `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt        time.Time         `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
	Token            []Token           `gorm:"foreignKey:user_id;references:id" json:"-"`
	UserRole         *UserRole         `gorm:"foreignKey:role_id;references:id" json:"user_role,omitempty"`
	UserSubscription *UserSubscription `gorm:"foreignKey:user_id;references:id" json:"user_subscription,omitempty"`
}

func (user *User) BeforeCreate(_ *gorm.DB) error {
	user.ID = uuid.New()
	return nil
}

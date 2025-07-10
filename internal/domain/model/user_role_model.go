package model

import (
	"time"
)

type UserRole struct {
	ID       uint   `gorm:"primaryKey"`
	RoleName string `gorm:"not null;unique"`

	Permission []RolePermissions `gorm:"many2many:user_role_permissions;"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

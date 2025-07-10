package model

import (
	"time"
)

type RolePermissions struct {
	ID             uint   `gorm:"primaryKey"`
	PermissionName string `gorm:"not null;unique"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

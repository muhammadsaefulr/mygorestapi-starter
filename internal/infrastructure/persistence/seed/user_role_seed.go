package seed

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// seed

func SeedUserRoles(db *gorm.DB) error {
	roles := []model.UserRole{
		{RoleName: "user"},
		{RoleName: "admin"},
	}

	for _, role := range roles {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error; err != nil {
			return err
		}
	}

	return nil
}

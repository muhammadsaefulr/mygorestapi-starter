package seed

import (
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedUserBadges(db *gorm.DB) error {
	badges := []model.UserBadge{
		{
			BadgeName: "Premium",
			IconURL:   "",
			Color:     "#FFD700",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, badge := range badges {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&badge).Error; err != nil {
			return err
		}
	}

	return nil
}

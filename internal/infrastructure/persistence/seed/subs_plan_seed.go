package seed

import (
	"time"

	"github.com/lib/pq"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedSubscriptionPlans(db *gorm.DB) error {
	plans := []model.SubscriptionPlan{
		{
			PlanName: "VIP Weekly",
			Duration: 168, // In Hours // Weekly
			Price:    10000,
			Benefits: pq.StringArray{
				"Tanpa iklan",
				"Badge premium",
				"Unlimited Poin Key",
				"Resolusi 1080p",
				"Komentar bisa lebih dari 3x per episode",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			PlanName: "VIP Monthly",
			Duration: 720, // In Hours // Monthly
			Price:    15000,
			Benefits: pq.StringArray{
				"Tanpa iklan",
				"Badge premium",
				"Unlimited Poin Key",
				"Resolusi 1080p",
				"Komentar bisa lebih dari 3x per episode",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, plan := range plans {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&plan).Error; err != nil {
			return err
		}
	}

	return nil
}

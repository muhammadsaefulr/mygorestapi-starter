package seed

import (
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedUsers(db *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []model.User{
		{
			Name:      "Fake Name",
			Email:     "fake@example.com",
			Password:  string(hashedPassword),
			RoleId:    1,
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Admin",
			Email:     "admin@dev.com",
			Password:  string(hashedPassword),
			RoleId:    2,
			Role:      "admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

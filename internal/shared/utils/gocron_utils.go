package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

func StartVIPCronJob(db *gorm.DB) {
	s := gocron.NewScheduler(time.Local)

	s.Every(1).Day().At("02:00").Do(func() {
		DeleteExpiredVIP(db)
	})

	log.Printf("Starting VIP cron job scheduler successfully")

	s.StartAsync()
}

func DeleteExpiredVIP(db *gorm.DB) {
	result := db.Where("end_date < ?", time.Now()).Delete(&model.UserSubscription{})
	fmt.Printf("[%s] Deleted expired VIP: %d\n", time.Now().Format(time.RFC3339), result.RowsAffected)
}

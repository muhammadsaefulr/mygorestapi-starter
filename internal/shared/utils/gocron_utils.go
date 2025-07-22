package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
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
	var expiredSubs []model.UserSubscription

	if err := db.Where("end_date < ?", time.Now()).Find(&expiredSubs).Error; err != nil {
		fmt.Println("Error fetching expired VIPs:", err)
		return
	}

	var userIDs []uuid.UUID
	for _, sub := range expiredSubs {
		userIDs = append(userIDs, sub.UserID)
	}

	if len(userIDs) > 0 {
		res := db.Where("user_id IN ?", userIDs).Delete(&model.UserBadgeInfo{})
		fmt.Printf("[%s] Deleted user badges: %d\n", time.Now().Format(time.RFC3339), res.RowsAffected)
	}

	res := db.Where("end_date < ?", time.Now()).Delete(&model.UserSubscription{})
	fmt.Printf("[%s] Deleted expired VIPs: %d\n", time.Now().Format(time.RFC3339), res.RowsAffected)
}

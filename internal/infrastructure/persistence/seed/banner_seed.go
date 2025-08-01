package seed

import (
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedBannerApp(db *gorm.DB) error {
	banners := []model.BannerApp{
		{
			Title:      "Mizu Zokusei no Mahoutsukai Subtitle Indonesia",
			ImageUrl:   "https://otakudesu.best/wp-content/uploads/2025/07/Mizu-Zokusei-no-Mahoutsukai.jpg",
			BannerType: "anime",
			DetailURL:  "/otakudesu/detail/mizu-zokusei-mahoutsukai-sub-indo",
			UpdatedBy:  "Admin",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Title:      "Dr. Stone: Science Future Part 2 Subtitle Indonesia",
			ImageUrl:   "https://otakudesu.best/wp-content/uploads/2025/07/148849.jpg",
			BannerType: "anime",
			DetailURL:  "/otakudesu/detail/ds-future-part2-sub-indo",
			UpdatedBy:  "Admin",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Title:      "One Piece Subtitle Indonesia",
			ImageUrl:   "https://otakudesu.best/wp-content/uploads/2021/05/One-Piece-Sub-Indo.jpg",
			BannerType: "anime",
			DetailURL:  "/otakudesu/detail/1piece-sub-indo",
			UpdatedBy:  "Admin",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Title:      "When Life Gives You Tangerines",
			ImageUrl:   "https://i.mydramalist.com/5v8b2y_4c.jpg?v=1",
			BannerType: "kdrama",
			DetailURL:  "/movie/details/when-life-gives-you-tangerines",
			UpdatedBy:  "Admin",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Title:      "Weak Hero Class 1",
			ImageUrl:   "https://i.mydramalist.com/pq2lr_4c.jpg?v=1",
			BannerType: "kdrama",
			DetailURL:  "/movie/details/weak-hero-class-1",
			UpdatedBy:  "Admin",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Title:      "Move to Heaven",
			ImageUrl:   "https://i.mydramalist.com/Rle36_4c.jpg?v=1",
			BannerType: "kdrama",
			DetailURL:  "/movie/details/move-to-heaven",
			UpdatedBy:  "Admin",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	for _, banner := range banners {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&banner).Error; err != nil {
			return err
		}
	}

	return nil
}

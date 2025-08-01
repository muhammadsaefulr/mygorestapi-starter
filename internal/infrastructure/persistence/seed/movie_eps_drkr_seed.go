package seed

import (
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedDrakorEpisodes(db *gorm.DB) error {
	movieIDs := []string{
		"strong-woman-do-bong-soon",
		"business-proposal",
		"descendants-of-the-sun",
		"the-trauma-code-heroes-on-call",
		"alchemy-of-souls",
		"moving",
		"twinkling-watermelon",
		"when-life-gives-you-tangerines",
		"weak-hero-class-1",
		"reply-1988",
		"flower-of-evil",
		"hospital-playlist-season-2",
		"move-to-heaven",
	}

	resolutions := []string{"360p", "480p", "720p", "1080p"}
	videoURL := "https://dev.msaepul.my.id/minio/drama/testing/testing.mp4"
	sourceBy := "Admin Server"

	for _, movieID := range movieIDs {
		episodeID := movieID + "-eps-01"

		for _, res := range resolutions {
			episode := model.MovieEpisode{
				MovieEpsID: episodeID,
				MovieId:    movieID,
				Resolution: res,
				VideoURL:   videoURL,
				Title:      "upload",
				SourceBy:   sourceBy,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&episode).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

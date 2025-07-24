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

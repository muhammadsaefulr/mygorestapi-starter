package seed

import (
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedMovieEpisodes(db *gorm.DB) error {
	movieIDs := []string{
		"the-fantastic-four-first-steps",
		"a-normal-woman",
		"lilo-and-stitch",
		"happy-gilmore-2",
		"how-to-train-your-dragon",
		"superman",
		"materialists",
		"demon-slayer-kimetsu-no-yaiba-infinity-castle",
		"xxx",
		"dangerous-animals",
		"m3gan-2",
		"man-with-no-past",
		"karate-kid-legends",
		"bride-hard",
		"ice-road-vengeance",
	}

	resolutions := []string{"360p", "480p", "720p", "1080p"}
	videoURL := "https://dev.msaepul.my.id/minio/movie/testing/testing.mp4"
	sourceBy := "Admin Server"

	for _, movieID := range movieIDs {
		episodeID := movieID + "-movie"

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

package seed

import (
	"time"

	"github.com/lib/pq"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedDrakor(db *gorm.DB) error {
	drakorSeed := []model.MovieDetails{
		{
			MovieID:      "strong-woman-do-bong-soon",
			MovieType:    "kdrama",
			ThumbnailURL: "https://image.tmdb.org/t/p/original/your_poster_path.jpg",
			Title:        "Strong Woman Do Bong Soon",
			Rating:       "8.7",
			Producer:     "JTBC",
			Status:       "Finished Airing",
			Studio:       "Drama House, JS Pictures",
			ReleaseDate:  "2017-02-24",
			Synopsis:     "Do Bong‑soon was born with superhuman strength. Recruited as bodyguard to a gaming‑chaebol heir, she uses her power to fight crime while juggling a love triangle with her childhood friend and her boss.",
			Genres:       pq.StringArray{"Fantasy", "Action", "Romance", "Comedy"},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			MovieID:      "business-proposal",
			MovieType:    "kdrama",
			ThumbnailURL: "https://hubpages.com/assets/resized/business-proposal-review.jpg",
			Title:        "Business Proposal",
			Rating:       "8.1",
			Producer:     "SBS TV",
			Status:       "Finished Airing",
			Studio:       "StudioS, Kakao Entertainment, Kross Pictures",
			ReleaseDate:  "2022-02-28",
			Synopsis:     "A romantic comedy where Shin Ha‑ri goes on a blind date disguised as her friend and ends up dating her CEO, Kang Tae‑moo.",
			Genres:       pq.StringArray{"Romantic Comedy", "Workplace", "Comedy"},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			MovieID:      "descendants-of-the-sun",
			MovieType:    "kdrama",
			ThumbnailURL: "https://redbubble.com/path-to-poster.jpg",
			Title:        "Descendants of the Sun",
			Rating:       "8.2",
			Producer:     "KBS2",
			Status:       "Finished Airing",
			Studio:       "KBS Drama Production, Next Entertainment World, Barunson Inc.",
			ReleaseDate:  "2016-02-24",
			Synopsis:     "A soldier Yoo Si‑jin and a surgeon Kang Mo‑yeon fall in love while facing dangers in a war‑like setting.",
			Genres:       pq.StringArray{"Romance", "Melodrama", "Action"},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, drakor := range drakorSeed {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&drakor).Error; err != nil {
			return err
		}
	}

	return nil
}

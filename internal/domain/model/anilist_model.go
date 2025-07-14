package model

import (
	"time"
)

type Anilist struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AniListMedia struct {
	ID    int `json:"id"`
	Title struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
		Native  string `json:"native"`
	} `json:"title"`
	CoverImage struct {
		Large string `json:"large"`
	} `json:"coverImage"`
	AverageScore int      `json:"averageScore"`
	Genres       []string `json:"genres"`
	Status       string   `json:"status"`
	Episodes     int      `json:"episodes"`
	NextAiring   struct {
		AiringAt int64 `json:"airingAt"`
		Episode  int   `json:"episode"`
	} `json:"nextAiringEpisode"`
	Studios struct {
		Nodes []struct {
			Name string `json:"name"`
		} `json:"nodes"`
	} `json:"studios"`
	StartDate struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	} `json:"startDate"`
	Description string `json:"description"`
}

type AniListResponse struct {
	Data struct {
		Page struct {
			PageInfo struct {
				PerPage     int  `json:"perPage"`
				CurrentPage int  `json:"currentPage"`
				LastPage    int  `json:"lastPage"`
				HasNextPage bool `json:"hasNextPage"`
			} `json:"pageInfo"`
			Media []AniListMedia `json:"media"`
		} `json:"Page"`
	} `json:"data"`
}

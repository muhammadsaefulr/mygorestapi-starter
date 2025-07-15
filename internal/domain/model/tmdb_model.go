package model

import (
	"time"
)

type Tmdb struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TMDbResponse struct {
	Page         int          `json:"page"`
	Results      []TMDbResult `json:"results"`
	TotalPages   int          `json:"total_pages"`
	TotalResults int          `json:"total_results"`
}

type TMDbResult struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"` // movie
	Name         string  `json:"name"`  // tv
	PosterPath   string  `json:"poster_path"`
	Overview     string  `json:"overview"`
	VoteAverage  float64 `json:"vote_average"`
	ReleaseDate  string  `json:"release_date"`   // movie
	FirstAirDate string  `json:"first_air_date"` // tv
}

type TMDbDetailResponse struct {
	Genres []struct {
		Name string `json:"name"`
	} `json:"genres"`

	ProductionCompanies []struct {
		Name string `json:"name"`
	} `json:"production_companies"`

	EpisodeRunTime   []int  `json:"episode_run_time"`
	NumberOfEpisodes int    `json:"number_of_episodes"`
	Status           string `json:"status"`
}

type TMDBGenreResponse struct {
	Genres []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
}

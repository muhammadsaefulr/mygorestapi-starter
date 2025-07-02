package response

import (
	"time"
)

type EpisodesResponse struct {
	MovieEpsId string `json:"movie_eps_id"`
	Title      string `json:"title"`
	VideoURL   string `json:"video_url"`
}

type MovieDetailOnlyResponse struct {
	MovieID      string    `json:"movie_id"`
	MovieType    string    `json:"movie_type"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Title        string    `json:"title"`
	Rating       string    `json:"rating"`
	Producer     string    `json:"producer"`
	Status       string    `json:"status"`
	TotalEps     string    `json:"total_eps"`
	Studio       string    `json:"studio"`
	ReleaseDate  string    `json:"release_date"`
	Synopsis     string    `json:"synopsis"`
	Genres       []string  `json:"genres"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MovieDetailsResponse struct {
	MovieDetail *MovieDetailOnlyResponse `json:"movie_details"`
	Episodes    []EpisodesResponse       `json:"episodes"`
	Rekomend    *MovieDetailOnlyResponse `json:"rekomend"`
}

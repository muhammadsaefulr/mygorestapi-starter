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
	MovieID      string     `json:"movie_id,omitempty"`
	MovieType    string     `json:"movie_type"`
	PathURL      string     `json:"path_url,omitempty"`
	ThumbnailURL string     `json:"thumbnail_url"`
	Title        string     `json:"title"`
	Rating       string     `json:"rating,omitempty"`
	Producer     string     `json:"producer,omitempty"`
	Status       string     `json:"status,omitempty"`
	TotalEps     string     `json:"total_eps,omitempty"`
	Studio       string     `json:"studio,omitempty"`
	ReleaseDate  string     `json:"release_date,omitempty"`
	Synopsis     string     `json:"synopsis,omitempty"`
	Genres       []string   `json:"genres"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type MovieDetailsResponse struct {
	MovieDetail *MovieDetailOnlyResponse `json:"movie_details"`
	Episodes    []EpisodesResponse       `json:"episodes"`
	Rekomend    *MovieDetailOnlyResponse `json:"rekomend"`
}

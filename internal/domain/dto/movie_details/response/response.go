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
	IDSource     string `json:"id_source,omitempty"`
	MovieID      string `json:"movie_id,omitempty"`
	MovieType    string `json:"movie_type"`
	PathURL      string `json:"path_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Title        string `json:"title"`
	Rating       string `json:"rating"`
	Producer     string `json:"producer"`
	Status       string `json:"status"`
	TotalEps     string `json:"total_eps"`
	Studio       string `json:"studio,omitempty"`
	ReleaseDate  string `json:"release_date"`
	UpdateDay    string `json:"update_day,omitempty"`
	Synopsis     string `json:"synopsis"`

	Genres    []string                   `json:"genres"`
	CreatedAt *time.Time                 `json:"created_at,omitempty"`
	UpdatedAt *time.Time                 `json:"updated_at,omitempty"`
	Rekomend  *[]MovieDetailOnlyResponse `json:"rekomend,omitempty"`
}

type MovieDetailsResponse struct {
	MovieDetail *MovieDetailOnlyResponse `json:"movie_details"`
	Episodes    []EpisodesResponse       `json:"episodes"`
	Rekomend    *MovieDetailOnlyResponse `json:"rekomend"`
}

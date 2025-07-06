package response

import (
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
)

type HistoryResponse struct {
	ID           uint   `json:"id"`
	UserId       string `json:"user_id"`
	MovieId      string `json:"movie_id"`
	MovieEpsId   string `json:"movie_eps_id"`
	PlaybackTime int    `json:"playback_time"`
	AnimeDetail  response.MovieDetailOnlyResponse
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

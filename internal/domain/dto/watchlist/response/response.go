package response

import "time"

type WatchlistResponse struct {
	ID        uint      `json:"id"`
	UserId    string    `json:"user_id"`
	MovieId   string    `json:"movie_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

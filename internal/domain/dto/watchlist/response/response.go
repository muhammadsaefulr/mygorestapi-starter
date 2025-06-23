package response

import "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"

type WatchlistResponse struct {
	ID          uint   `json:"id"`
	UserId      string `json:"user_id"`
	MovieId     string `json:"movie_id"`
	AnimeDetail model.AnimeDetail
}

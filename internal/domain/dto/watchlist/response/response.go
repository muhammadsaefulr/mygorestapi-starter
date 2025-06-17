package response

type WatchlistResponse struct {
	ID      uint   `json:"id"`
	UserId  string `json:"user_id"`
	MovieId string `json:"movie_id"`
}

package request

type CreateWatchlist struct {
	UserId  string `json:"-"`
	MovieId string `json:"movie_id" validate:"required"`
}

type UpdateWatchlist struct {
	ID      uint   `json:"-"`
	UserId  string `json:"-"`
	MovieId string `json:"movie_id" validate:"required"`
}

type QueryWatchlist struct {
	Search string `query:"search"`
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Sort   string `query:"sort"`
}

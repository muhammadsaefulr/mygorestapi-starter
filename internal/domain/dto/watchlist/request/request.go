package request

type CreateWatchlist struct {
	UserId  string `json:"user_id" validate:"required,uuid"`
	MovieId string `json:"movie_id" validate:"required,uuid"`
}

type UpdateWatchlist struct {
	ID      uint   `json:"-"`
	UserId  string `json:"user_id" validate:"required,uuid"`
	MovieId string `json:"movie_id" validate:"required,uuid"`
}

type QueryWatchlist struct {
	Search string `query:"search"`
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Sort   string `query:"sort"`
}

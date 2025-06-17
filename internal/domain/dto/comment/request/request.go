package request

type CreateComment struct {
	UserId  string `json:"-"`
	MovieId string `json:"movie_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type UpdateComment struct {
	ID      uint   `json:"-"`
	UserId  string `json:"-"`
	MovieId string `json:"-"`
	Content string `json:"content" validate:"required"`
}

type QueryComment struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

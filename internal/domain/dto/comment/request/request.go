package request

type CreateComment struct {
	Name string `json:"name"`
}

type UpdateComment struct {
	ID uint `json:"-"`
	Name string `json:"name"`
}

type QueryComment struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

package request

type CreateTmdb struct {
	Name string `json:"name"`
}

type UpdateTmdb struct {
	ID   uint   `json:"-"`
	Name string `json:"name"`
}

type QueryTmdb struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Category string `query:"category"` // contoh: "popular", "ongoing", "trending"
	Search   string `query:"search"`
	Type     string `query:"type"` // contoh: "movie", "kdrama"
}

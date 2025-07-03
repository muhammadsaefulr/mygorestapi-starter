package request

type CreateAnilist struct {
	Name string `json:"name"`
}

type UpdateAnilist struct {
	ID   uint   `json:"-"`
	Name string `json:"name"`
}

type QueryAnilist struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Sort     string `query:"sort"`
	Search   string `query:"search"`
	Category string `query:"category"` // contoh: "popular", "ongoing", "trending"
	Rekom    string `query:"rekom"`    // optional, contoh: "one piece"
}

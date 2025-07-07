package request

type CreateMdl struct {
	Name string `json:"name"`
}

type UpdateMdl struct {
	ID   uint   `json:"-"`
	Name string `json:"name"`
}

type QueryMdl struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Category string `query:"category"`
	Sort     string `query:"sort"`
}

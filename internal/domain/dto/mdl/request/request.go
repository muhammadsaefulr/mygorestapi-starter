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
	Search   string `query:"search"`
	Limit    int    `query:"limit"`
	Category string `query:"category"`
	Genre    string `query:"genre"`
	Sort     string `query:"sort"`
}

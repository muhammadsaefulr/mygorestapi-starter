package request

type CreateUserRole struct {
	Name       string `json:"name"`
	Permission []uint `json:"permission"`
}

type UpdateUserRole struct {
	ID         uint   `json:"-"`
	Name       string `json:"name"`
	Permission []uint `json:"permission"`
}

type QueryUserRole struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Sort   string `query:"sort"`
	Search string `query:"search"`
}

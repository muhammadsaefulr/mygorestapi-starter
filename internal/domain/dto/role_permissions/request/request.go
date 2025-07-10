package request

type CreateRolePermissions struct {
	Name string `json:"name"`
}

type UpdateRolePermissions struct {
	ID uint `json:"-"`
	Name string `json:"name"`
}

type QueryRolePermissions struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

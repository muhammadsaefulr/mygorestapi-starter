package request

type CreateDiscovery struct {
	Name string `json:"name"`
}

type UpdateDiscovery struct {
	ID   uint   `json:"-"`
	Name string `json:"name"`
}

type QueryDiscovery struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Category string `query:"category"`
	Search   string `query:"search"`
	Type     string `query:"type"`
}

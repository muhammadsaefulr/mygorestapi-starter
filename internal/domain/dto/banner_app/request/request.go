package request

type CreateBannerApp struct {
	Title      string `json:"title"`
	ImageUrl   string `json:"image_url"`
	BannerType string `json:"banner_type"` // e.g., "movie", "drakor", "anime"
	UpdatedBy  string `json:"-"`           // typically the User Name of the person creating the banner
	DetailURL  string `json:"detail_url"`
}

type UpdateBannerApp struct {
	ID         uint   `json:"-"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	ImageUrl   string `json:"image_url"`
	BannerType string `json:"banner_type"`
	UpdatedBy  string `json:"-"` // typically the User Name of the person updating the banner
	DetailURL  string `json:"detail_url"`
}

type QueryBannerApp struct {
	Page  int    `query:"page"`
	Type  string `query:"type"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

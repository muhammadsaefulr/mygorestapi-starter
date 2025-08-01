package request

type CreateBannerApp struct {
	Title     string `json:"title"`
	ImageUrl  string `json:"image_url"`
	DetailURL string `json:"detail_url"`
}

type UpdateBannerApp struct {
	ID        uint   `json:"-"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	ImageUrl  string `json:"image_url"`
	DetailURL string `json:"detail_url"`
}

type QueryBannerApp struct {
	Page  int    `query:"page"`
	Type  string `query:"type"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

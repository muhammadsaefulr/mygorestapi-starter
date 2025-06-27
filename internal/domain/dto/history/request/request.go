package request

type CreateHistory struct {
	UserId       string `json:"-"`
	MovieEpsId   string `json:"movie_eps_id"`
	PlaybackTime int    `json:"playback_time"`
}

type UpdateHistory struct {
	ID           uint   `json:"-"`
	UserId       string `json:"-"`
	MovieEpsId   string `json:"-"`
	PlaybackTime int    `json:"playback_time"`
}

type QueryHistory struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

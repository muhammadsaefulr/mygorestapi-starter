package response

type MovieEpisodesDetails struct {
	ID         uint   `json:"id,omitemty"`
	MovieEpsID string `json:"movie_eps_id"`
	MovieId    string `json:"movie_id"`
	Resolution string `json:"resolution"`
	VideoURL   string `json:"video_url"`
	Title      string `json:"title"`
}

type Sources struct {
	DataList   []SourcesData `json:"data_list"`
	Res        string        `json:"res"`
	MovieEpsId string        `json:"movie_eps_id"`
}

type SourcesData struct {
	Title    string `json:"title"`
	VideoURL string `json:"video_url"`
}

type MovieEpisodeResponses struct {
	Title        string    `json:"title"`
	ReleaseDate  string    `json:"release_date"`
	ThumbnailURL string    `json:"thumbnail_url"`
	CurrentEp    string    `json:"current_ep"`
	DetailUrl    string    `json:"detail_url"`
	NextEpUrl    string    `json:"next_ep_url"`
	Sources      []Sources `json:"sources"`
	Episodes     []SourcesData
}

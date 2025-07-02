package response

import "time"

type MovieEpisodesDetails struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

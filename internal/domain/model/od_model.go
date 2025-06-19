package model

type AnimeData struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	JudulPath    string `json:"judul_path"`
	ThumbnailURL string `json:"thumbnail_url"`
	LatestEp     string `json:"latest_ep"`
	UpdateAnime  string `json:"update_anime"`
}

type TrendingAnime struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	JudulPath    string `json:"judul_path"`
	ThumbnailURL string `json:"thumbnail_url"`
	LatestEp     string `json:"latest_ep"`
	UpdateAnime  string `json:"update_anime"`
}

type OngoingAnime struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	JudulPath    string `json:"judul_path"`
	ThumbnailURL string `json:"thumbnail_url"`
	Episode      string `json:"episode"`
	DaysUpdated  string `json:"days_updated"`
	UpdatedAt    string `json:"updated_at"`
}

type CompleteAnime struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	JudulPath    string `json:"judul_path"`
	ThumbnailURL string `json:"thumbnail_url"`
	LatestEp     string `json:"latest_ep"`
	Rating       string `json:"rating"`
	UpdatedAt    string `json:"updated_at"`
}

type GenreAnime struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Studio   string `json:"studio"`
	Episodes string `json:"episodes"`
	Rating   string `json:"rating"`
}

type GenreInfo struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type SearchResult struct {
	Title        string      `json:"title"`
	URL          string      `json:"url"`
	ThumbnailURL string      `json:"thumbnail_url"`
	Genres       []GenreInfo `json:"genres"`
	Status       string      `json:"status"`
	Rating       string      `json:"rating"`
}

// Anime Episode Types

type AnimeEpisode struct {
	Title    string `json:"title"`
	VideoURL string `json:"video_url"`
}

type AnimeDetail struct {
	ThumbnailURL string      `json:"thumbnail_url"`
	Title        string      `json:"title"`
	Rating       string      `json:"rating"`
	Producer     string      `json:"producer"`
	Status       string      `json:"status"`
	TotalEps     string      `json:"total_eps"`
	Duration     string      `json:"duration"`
	Studio       string      `json:"studio"`
	ReleaseDate  string      `json:"release_date"`
	Synopsis     string      `json:"synopsis"`
	Genres       []GenreInfo `json:"genres"`
}

type EpisodePageResult struct {
	AnimeDetail AnimeDetail    `json:"anime_detail"`
	AnimeEps    []AnimeEpisode `json:"episode"`
}

// Anime Video Source Data

type VideoSource struct {
	Res      string         `json:"res"`
	DataList []AnimeEpisode `json:"data_list"`
}

// type SourceLink struct {
// 	Title string `json:"title"`
// 	URL   string `json:"url"`
// }

type AnimeSourceData struct {
	Title       string         `json:"title"`
	ReleaseDate string         `json:"release_date"`
	CurrentEp   string         `json:"current_ep"`
	DownloadURL string         `json:"download_url"`
	NextEpURL   string         `json:"next_ep_url"`
	Sources     []VideoSource  `json:"sources"`
	Episodes    []AnimeEpisode `json:"episodes"`
}

package request

import "mime/multipart"

type CreateMovieEpisodes struct {
	MovieEpsID     string `json:"movie_eps_id" example:"wu-nao-monu-eps-1"`
	MovieId        string `json:"movie_id" example:"wu-nao-monu"`
	Title          string `json:"title" example:"upload"`
	Resolution     string `json:"resolution" example:"720p"`
	ContentUploads string `json:"video_url" example:"youtube"`
	SourceBy       string `json:"-"`
}

type UpdateMovieEpisodes struct {
	MovieEpsID     string `json:"-"`
	MovieId        string `json:"movie_id"`
	Title          string `json:"title"`
	Resolution     string `json:"resolution"`
	ContentUploads string `json:"url"`
}

type CreateMovieEpisodesUpload struct {
	MovieEpsID     string                `form:"movie_eps_id" validate:"required" example:"wu-nao-monu-eps-1"`
	MovieId        string                `form:"movie_id" validate:"required" example:"wu-nao-monu"`
	Title          string                `form:"title" validate:"required" example:"upload"`
	Resolution     string                `form:"resolution" validate:"required" example:"720p"`
	ContentUploads *multipart.FileHeader `form:"file_video" validate:"required"`
	SourceBy       string                `form:"-"`
}

type UpdateMovieEpisodesUpload struct {
	MovieEpsID     string                `form:"-"`
	MovieId        string                `form:"movie_id" example:"wu-nao-monu"`
	Title          string                `form:"title"  example:"upload"`
	Resolution     string                `form:"resolution"  example:"720p"`
	ContentUploads *multipart.FileHeader `form:"file_video" validate:"required"`
	SourceBy       string                `form:"-"` // tetap diisi dari backend
}

type QueryMovieEpisode struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Sort   string `query:"sort"`
	Search string `query:"search"`
}

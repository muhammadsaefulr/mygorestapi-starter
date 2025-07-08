package request

type CreateMovieDetails struct {
	MovieID      string   `json:"movie_id" example:"wu-nao-monu" validate:"required"`
	MovieType    string   `json:"movie_type" example:"anime" validate:"required"`
	ThumbnailURL string   `json:"thumbnail_url" example:"https://s4.anilist.co/file/anilistcdn/media/anime/cover/medium/bx141914-P1xQHMXN7q6z.png" validate:"required,url"`
	Title        string   `json:"title" example:"Wu Nao Monü" validate:"required"`
	Rating       string   `json:"rating" example:"5.9"`
	Producer     string   `json:"producer" example:"Agate"`
	Status       string   `json:"status" example:"complete"`
	Studio       string   `json:"studio" example:"Agate"`
	ReleaseDate  string   `json:"release_date" example:"2023-05-01"`
	Synopsis     string   `json:"synopsis" example:"Agate is a demon girl cursed with an eternal life. To forget the past, she throws half of her head into a deep valley and runs away. Unexpectedly, the tears flood the valley floor and form into a lake, triggering a flood in hell. Aloys, the Prince of Ghost comes to the Human World to find out the truth and finally finds the demon girl, starting a story of life and love."`
	Genres       []string `json:"genres" example:"demon,fantasy,supernatural"  validate:"required"`
}

type UpdateMovieDetails struct {
	MovieID      string   `json:"-" example:"wu-nao-monu"`
	MovieType    string   `json:"movie_type" example:"anime" validate:"required"`
	ThumbnailURL string   `json:"thumbnail_url" example:"https://s4.anilist.co/file/anilistcdn/media/anime/cover/medium/bx141914-P1xQHMXN7q6z.png" validate:"required,url"`
	Title        string   `json:"title" example:"Wu Nao Monü" validate:"required"`
	Rating       string   `json:"rating" example:"5.9"`
	Producer     string   `json:"producer" example:"Agate"`
	Status       string   `json:"status" example:"complete"`
	Studio       string   `json:"studio" example:"Agate"`
	ReleaseDate  string   `json:"release_date" example:"2023-05-01"`
	Synopsis     string   `json:"synopsis" example:"Agate is a demon girl cursed with an eternal life. To forget the past, she throws half of her head into a deep valley and runs away. Unexpectedly, the tears flood the valley floor and form into a lake, triggering a flood in hell. Aloys, the Prince of Ghost comes to the Human World to find out the truth and finally finds the demon girl, starting a story of life and love."`
	Genres       []string `json:"genres" example:"demon,fantasy,supernatural"  validate:"required"`
}

type QueryMovieDetails struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Type   string `query:"type"`
	Sort   string `query:"sort"`
	Search string `query:"search"`
}

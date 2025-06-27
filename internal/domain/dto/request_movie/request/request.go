package request

import "time"

type CreateRequestMovie struct {
	Title         string   `json:"title" example:"Jujutsu Kaisen Season 2"`
	UserIdRequest string   `json:"-"`
	TypeMovie     string   `json:"type_movie" example:"Anime"`
	Genre         []string `json:"genre" example:"Action,Shounen"`
	Description   string   `json:"description" example:"Yuuji Itadori joins Jujutsu Tech to exorcise curses."`
	StatusMovie   string   `json:"status_movie" example:"Ongoing"`
	StatusRequest string   `json:"status_request" example:"Pending"`
}

type UpdateRequestMovie struct {
	ID            uint      `json:"-"`
	Title         string    `json:"title" example:"Jujutsu Kaisen Season 2"`
	TypeMovie     string    `json:"type_movie" example:"TV"`
	Genre         []string  `json:"genre" example:"Action,Shounen"`
	Description   string    `json:"description" example:"Yuuji Itadori joins Jujutsu Tech to exorcise curses."`
	StatusMovie   string    `json:"status_movie" example:"Ongoing"`
	StatusRequest string    `json:"status_request" example:"Pending"`
	UpdatedAt     time.Time `json:"-"`
}

type QueryRequestMovie struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Sort   string `query:"sort"`
	Search string
}

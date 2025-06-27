package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID   uuid.UUID `json:"id" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	Name string    `json:"name" example:"Saeful"`
	Role string    `json:"role" example:"user"`
}

type RequestMovieResponse struct {
	ID            uint         `json:"id" example:"1"`
	UserIdRequest uuid.UUID    `json:"-"`
	Title         string       `json:"title" example:"Jujutsu Kaisen Season 2"`
	TypeMovie     string       `json:"type_movie" example:"TV"`
	Genre         []string     `json:"genre" example:"[\"Action\", \"Shounen\"]"`
	Description   string       `json:"description" example:"Yuuji Itadori joins Jujutsu Tech to exorcise curses."`
	StatusMovie   string       `json:"status_movie" example:"Ongoing"`
	StatusRequest string       `json:"status_request" example:"Pending"`
	RequestedBy   UserResponse `json:"requested_by"`
	CreatedAt     time.Time    `json:"created_at" example:"2025-06-27T15:04:05Z"`
	UpdatedAt     time.Time    `json:"updated_at" example:"2025-06-27T15:04:05Z"`
}

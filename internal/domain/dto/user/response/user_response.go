package response

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type CreateUserResponse struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	IsEmailVerified bool   `json:"is_email_verified"`
}

type GetUsersResponse struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Email           string          `json:"email"`
	Role            string          `json:"role"`
	Roles           *model.UserRole `json:"roles"`
	IsEmailVerified bool            `json:"is_email_verified"`
}

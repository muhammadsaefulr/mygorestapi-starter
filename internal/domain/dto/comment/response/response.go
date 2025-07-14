package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/user/response"
)

type CommentResponse struct {
	ID           uint                       `json:"id"`
	UserId       uuid.UUID                  `json:"user_id"`
	MovieId      string                     `json:"movie_id"`
	Content      string                     `json:"content"`
	CreatedAt    time.Time                  `json:"created_at"`
	UserDetal    *response.GetUsersResponse `json:"user_detail,omitempty"`
	CommentReply []CommentResponse          `json:"comment_reply"`
	Likes        int                        `json:"likes"`
}

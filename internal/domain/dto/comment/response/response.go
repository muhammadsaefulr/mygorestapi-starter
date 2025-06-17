package response

import (
	"time"

	"github.com/google/uuid"
)

type CommentResponse struct {
	// contoh field, sesuaikan dengan kebutuhan
	ID        uint      `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	MovieId   string    `json:"movie_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

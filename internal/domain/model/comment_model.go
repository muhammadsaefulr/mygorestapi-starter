package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID      uint      `gorm:"primaryKey;not_null" json:"id"`
	UserId  uuid.UUID `gorm:"not null" json:"user_id"`
	MovieId string    `gorm:"not null" json:"movie_id"`
	Content string    `gorm:"not null" json:"content"`

	ParentId *uint     `json:"parent_id"`
	Replies  []Comment `gorm:"foreignKey:ParentId;references:ID" json:"replies,omitempty"`

	Likes []CommentLike `gorm:"foreignKey:CommentID;references:ID" json:"likes,omitempty"`

	UserDetail User      `gorm:"foreignKey:UserId;references:ID" json:"user_detail"`
	CreatedAt  time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`
}

type CommentLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserId    uuid.UUID `gorm:"not null;uniqueIndex:idx_user_comment" json:"user_id"`
	CommentID uint      `gorm:"not null;uniqueIndex:idx_user_comment" json:"comment_id"`

	User      User      `gorm:"foreignKey:UserId;references:ID" json:"user"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
}

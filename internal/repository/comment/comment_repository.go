package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetCommentByID(ctx context.Context, id uint) (*model.Comment, error)
	GetCommentByMovieID(ctx context.Context, movieID string) ([]model.Comment, error)
	UpdateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	DeleteComment(ctx context.Context, id uint) error
	LikeComment(ctx context.Context, dataCmnt model.CommentLike) error
	HasUserLiked(ctx context.Context, commentID uint, userID string) (bool, error)
	DislikeComment(ctx context.Context, id uint, userId string) error
}

package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (c *commentRepository) CreateComment(ctx context.Context, comment *model.Comment) error {
	return c.db.WithContext(ctx).Create(comment).Error
}

func (c *commentRepository) DeleteComment(ctx context.Context, id uint) error {
	return c.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Comment{}).Error
}

func (c *commentRepository) GetCommentByMovieID(ctx context.Context, movieID string) ([]model.Comment, error) {
	var comments []model.Comment
	if err := c.db.WithContext(ctx).Where("movie_id = ?", movieID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *commentRepository) GetCommentByID(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := c.db.WithContext(ctx).First(&comment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *commentRepository) UpdateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	if err := c.db.WithContext(ctx).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

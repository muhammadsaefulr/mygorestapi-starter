package repository

import (
	"context"
	"time"

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
	if err := c.db.WithContext(ctx).
		Preload("UserDetail").
		Where("movie_id = ?", movieID).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *commentRepository) GetCommentByID(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := c.db.WithContext(ctx).
		Preload("UserDetail").
		First(&comment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *commentRepository) UpdateComment(ctx context.Context, updated *model.Comment) (*model.Comment, error) {
	var comment model.Comment
	if err := c.db.WithContext(ctx).First(&comment, updated.ID).Error; err != nil {
		return nil, err
	}

	comment.Content = updated.Content
	comment.UpdatedAt = time.Now()

	if err := c.db.WithContext(ctx).Save(&comment).Error; err != nil {
		return nil, err
	}

	if err := c.db.WithContext(ctx).
		Preload("UserDetail").
		First(&comment, comment.ID).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

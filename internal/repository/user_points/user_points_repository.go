package repository

import (
	"context"

	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type UserPointsRepo interface {
	GetByUserID(ctx context.Context, id string) (*model.UserPoints, error)
	Update(ctx context.Context, data *model.UserPoints) error
	Create(ctx context.Context, data *model.UserPoints) error
	Delete(ctx context.Context, id string) error
}

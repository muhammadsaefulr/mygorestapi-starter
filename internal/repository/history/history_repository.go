package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/history/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type HistoryRepo interface {
	GetAllByUserId(ctx context.Context, UserId string, param *request.QueryHistory) ([]model.History, int64, error)
	GetByID(ctx context.Context, id uint) (*model.History, error)
	Create(ctx context.Context, data *model.History) error
	Update(ctx context.Context, data *model.History) error
	Delete(ctx context.Context, id uint) error
}

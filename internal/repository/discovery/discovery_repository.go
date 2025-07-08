package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/discovery/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type DiscoveryRepo interface {
	GetAll(ctx context.Context, param *request.QueryDiscovery) ([]model.Discovery, int64, error)
	GetByID(ctx context.Context, id uint) (*model.Discovery, error)
	Create(ctx context.Context, data *model.Discovery) error
	Update(ctx context.Context, data *model.Discovery) error
	Delete(ctx context.Context, id uint) error
}

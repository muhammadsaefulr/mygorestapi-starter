package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/request_vip/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type RequestVipRepo interface {
	GetAll(ctx context.Context, param *request.QueryRequestVip) ([]model.RequestVip, int64, error)
	GetByID(ctx context.Context, id uint) (*model.RequestVip, error)
	Create(ctx context.Context, data *model.RequestVip) error
	Update(ctx context.Context, data *model.RequestVip) error
	Delete(ctx context.Context, id uint) error
}

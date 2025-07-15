package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/banner_app/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type BannerAppRepo interface {
	GetAll(ctx context.Context, param *request.QueryBannerApp) ([]model.BannerApp, int64, error)
	GetByID(ctx context.Context, id uint) (*model.BannerApp, error)
	Create(ctx context.Context, data *model.BannerApp) error
	Update(ctx context.Context, data *model.BannerApp) error
	Delete(ctx context.Context, id uint) error
}

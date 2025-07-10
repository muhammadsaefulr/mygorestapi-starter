package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/report_error/request"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type ReportErrorRepo interface {
	GetAll(ctx context.Context, param *request.QueryReportError) ([]model.ReportError, int64, error)
	GetByID(ctx context.Context, id string) (*model.ReportError, error)
	Create(ctx context.Context, data *model.ReportError) error
	Update(ctx context.Context, data *model.ReportError) error
	Delete(ctx context.Context, id string) error
}

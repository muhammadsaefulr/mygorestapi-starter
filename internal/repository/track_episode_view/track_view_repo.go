package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

type TrackEpisodeViewRepository interface {
	GetAll(ctx context.Context) ([]TrackEpisodeViewSummary, error)
	Create(ctx context.Context, trackEpisodeView model.TrackEpisodeView) error
	GetByEpisodeId(ctx context.Context, episodeId string) (model.TrackEpisodeView, error)
}

package repository

import (
	"context"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"gorm.io/gorm"
)

type TrackEpisodeViewRepositoryImpl struct {
	db *gorm.DB
}

type TrackEpisodeViewSummary struct {
	MovieDetailUrl string `json:"movie_detail_url"`
	EpisodeId      string `json:"episode_id"`
	ViewCount      int    `json:"view_count"`
}

type TrackEpsParam struct {
	Page      int
	Limit     int
	MovieType string
}

func NewTrackEpisodeViewRepository(db *gorm.DB) TrackEpisodeViewRepository {
	return &TrackEpisodeViewRepositoryImpl{db: db}
}

func (t *TrackEpisodeViewRepositoryImpl) Create(ctx context.Context, trackEpisodeView model.TrackEpisodeView) error {
	return t.db.WithContext(ctx).Create(&trackEpisodeView).Error
}

func (t *TrackEpisodeViewRepositoryImpl) GetAll(ctx context.Context, param *TrackEpsParam) ([]TrackEpisodeViewSummary, error) {
	var results []TrackEpisodeViewSummary

	offset := (param.Page - 1) * param.Limit

	subQuery := t.db.
		Table("track_episode_views").
		Select("DISTINCT ON (movie_detail_url) movie_detail_url, episode_id, COUNT(*) as view_count").
		Group("movie_detail_url, episode_id").
		Order("movie_detail_url, COUNT(*) DESC")

	err := t.db.Table("(?) as ranked", subQuery).
		Order("view_count DESC").
		Limit(param.Limit).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (t *TrackEpisodeViewRepositoryImpl) GetByEpisodeId(ctx context.Context, episodeId string) (model.TrackEpisodeView, error) {
	var trackEpisodeView model.TrackEpisodeView
	if err := t.db.WithContext(ctx).Where("movie_detail_url = ?", episodeId).First(&trackEpisodeView).Error; err != nil {
		return model.TrackEpisodeView{}, err
	}
	return trackEpisodeView, nil
}

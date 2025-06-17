package convert_types

import (
	"github.com/google/uuid"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/watchlist/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateWatchlistToModel(req *request.CreateWatchlist) *model.Watchlist {
	return &model.Watchlist{
		UserId:        uuid.MustParse(req.UserId),
		MovieId:       req.MovieId,
		ThumbImageUrl: req.ThumbImageUrl,
	}
}

func UpdateWatchlistToModel(req *request.UpdateWatchlist) *model.Watchlist {
	return &model.Watchlist{
		ID:            req.ID,
		UserId:        uuid.MustParse(req.UserId),
		MovieId:       req.MovieId,
		ThumbImageUrl: req.ThumbImageUrl,
	}
}

func WathclistResponseToModel(watchlist *response.WatchlistResponse) *model.Watchlist {
	return &model.Watchlist{
		ID:      watchlist.ID,
		UserId:  uuid.MustParse(watchlist.UserId),
		MovieId: watchlist.MovieId,
	}
}

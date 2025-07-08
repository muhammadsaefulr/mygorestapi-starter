package convert_types

import (
	requestAn "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/anilist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/discovery/request"
	requestMdl "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/mdl/request"
	requestMvDtl "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/request"
	responseMvDtl "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	requestTm "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/tmdb/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func MapToAnilistQuery(q *request.QueryDiscovery) *requestAn.QueryAnilist {
	return &requestAn.QueryAnilist{
		Page:     q.Page,
		Limit:    q.Limit,
		Search:   q.Search,
		Category: q.Category,
	}
}

func MapToTmdbQuery(q *request.QueryDiscovery) *requestTm.QueryTmdb {
	return &requestTm.QueryTmdb{
		Page:     q.Page,
		Limit:    q.Limit,
		Search:   q.Search,
		Type:     q.Type,
		Category: q.Category,
	}
}

func MapToMdlQuery(q *request.QueryDiscovery) *requestMdl.QueryMdl {
	return &requestMdl.QueryMdl{
		Page:     q.Page,
		Limit:    q.Limit,
		Category: q.Category,
	}
}

func MapToMovieDtQuery(q *request.QueryDiscovery) *requestMvDtl.QueryMovieDetails {
	return &requestMvDtl.QueryMovieDetails{
		Page:   q.Page,
		Limit:  q.Limit,
		Search: q.Search,
	}
}

func ConvertMvDetailToOnlyResp(data []model.MovieDetails) []responseMvDtl.MovieDetailOnlyResponse {
	var results []responseMvDtl.MovieDetailOnlyResponse
	for _, d := range data {
		results = append(results, responseMvDtl.MovieDetailOnlyResponse{
			MovieID:      d.MovieID,
			MovieType:    d.MovieType,
			ThumbnailURL: d.ThumbnailURL,
			Title:        d.Title,
			Rating:       d.Rating,
			Producer:     d.Producer,
			Status:       d.Status,
			Studio:       d.Studio,
			ReleaseDate:  d.ReleaseDate,
			Synopsis:     d.Synopsis,
			Genres:       d.Genres,
			CreatedAt:    &d.CreatedAt,
			UpdatedAt:    &d.UpdatedAt,
		})
	}
	return results
}

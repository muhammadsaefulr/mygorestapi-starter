package convert_types

import (
	requestAn "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/anilist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/discovery/request"
	requestMdl "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/mdl/request"
	requestTm "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/tmdb/request"
)

func MapToAnilistQuery(q *request.QueryDiscovery) *requestAn.QueryAnilist {
	return &requestAn.QueryAnilist{
		Page:     q.Page,
		Limit:    q.Limit,
		Sort:     q.Sort,
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
		Sort:     q.Sort,
	}
}

package convert_types

import (
	"fmt"

	"strconv"
	"strings"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_episode/response"
	model "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func CreateMovieEpisodesToModel(req *request.CreateMovieEpisodes) *model.MovieEpisode {
	return &model.MovieEpisode{
		MovieEpsID: req.MovieEpsID,
		MovieId:    req.MovieId,
		Title:      req.Title,
		SourceBy:   req.SourceBy,
		Resolution: req.Resolution,
		VideoURL:   req.ContentUploads,
	}
}

func UpdateMovieEpisodesToModel(req *request.UpdateMovieEpisodes) *model.MovieEpisode {
	return &model.MovieEpisode{
		MovieEpsID: req.MovieEpsID,
		MovieId:    req.MovieId,
		Resolution: req.Resolution,
		VideoURL:   req.ContentUploads,
	}
}

func MovieEpisodeToResp(movieEpsList []model.MovieEpisode, movieDetails model.MovieDetails, movieEpsId string) response.MovieEpisodeResponses {
	groupedSources := make(map[string][]response.SourcesData)
	episodeSeen := make(map[string]bool)
	episodeList := make([]response.SourcesData, 0)

	for _, ep := range movieEpsList {
		source := response.SourcesData{
			Title:    ep.Title,
			VideoURL: ep.VideoURL,
		}
		groupedSources[ep.Resolution] = append(groupedSources[ep.Resolution], source)

		if !episodeSeen[ep.MovieEpsID] {
			episodeSeen[ep.MovieEpsID] = true
			episodeList = append(episodeList, response.SourcesData{
				Title:    strings.Title(strings.ReplaceAll(strings.ReplaceAll(ep.MovieEpsID, "-eps-", "-episode-"), "-", " ")),
				VideoURL: fmt.Sprintf("/movie/episodes/%s/%s", movieDetails.MovieID, ep.MovieEpsID),
			})
		}
	}

	sources := make([]response.Sources, 0, len(groupedSources))
	for res, dataList := range groupedSources {
		sources = append(sources, response.Sources{
			Res:        res,
			DataList:   dataList,
			MovieEpsId: movieEpsId,
		})
	}

	parts := strings.Split(movieEpsId, "-")
	currNum, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		currNum = -1
	}

	nextEpUrl := "unknown"
	expectedSuffix := fmt.Sprintf("eps-%d", currNum+1)
	for _, ep := range movieEpsList {
		if strings.HasSuffix(ep.MovieEpsID, expectedSuffix) {
			nextEpUrl = fmt.Sprintf("/movie/episodes/%s/%s", movieDetails.MovieID, ep.MovieEpsID)
			break
		}
	}

	return response.MovieEpisodeResponses{
		Title:        movieDetails.Title,
		ReleaseDate:  movieDetails.ReleaseDate,
		ThumbnailURL: movieDetails.ThumbnailURL,
		CurrentEp:    fmt.Sprintf("Episode %d", currNum),
		DetailUrl:    fmt.Sprintf("/movie/detail/%s", movieDetails.MovieID),
		NextEpUrl:    nextEpUrl,
		Sources:      sources,
		Episodes:     episodeList,
	}
}

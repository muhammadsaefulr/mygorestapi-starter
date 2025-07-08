package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/tmdb/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

var tmdbApiKey = "7a69a6a33b39eb0ba8503c9fe32eb132"

func FetchTMDbMedia(category, queryTitle, mediaType string, param *request.QueryTmdb) ([]response.MovieDetailOnlyResponse, error) {
	baseURL := "https://api.themoviedb.org/3"
	var endpoint string

	limit := param.Limit
	if limit == 0 {
		limit = 20
	}

	switch category {
	case "popular":
		endpoint = fmt.Sprintf("%s/%s/popular", baseURL, mediaType)
	case "trending":
		endpoint = fmt.Sprintf("%s/trending/%s/day", baseURL, mediaType)
	case "ongoing":
		if mediaType != "tv" {
			return nil, fmt.Errorf("ongoing hanya tersedia untuk tv show (kdrama)")
		}
		endpoint = fmt.Sprintf("%s/tv/on_the_air", baseURL)
	case "search":
		if queryTitle == "" {
			return nil, fmt.Errorf("query title tidak boleh kosong untuk pencarian")
		}
		endpoint = fmt.Sprintf("%s/search/%s?query=%s", baseURL, mediaType, url.QueryEscape(queryTitle))
	default:
		return nil, fmt.Errorf("invalid category")
	}

	endpoint += fmt.Sprintf("?api_key=%s&page=%d", tmdbApiKey, param.Page)
	if mediaType == "tv" {
		endpoint += "&with_origin_country=KR"
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		// fmt.Println(">>> [ERROR] Failed fetch main list:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result model.TMDbResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println(">>> [ERROR] Decode error:", err)
		return nil, err
	}

	if len(result.Results) == 0 {
		fmt.Println(">>> [WARNING] Empty result from TMDb")
		return nil, nil
	}

	if len(result.Results) > limit {
		result.Results = result.Results[:limit]
	}

	var (
		movies    = make([]response.MovieDetailOnlyResponse, len(result.Results))
		wg        sync.WaitGroup
		semaphore = make(chan struct{}, 10)
	)

	for i, item := range result.Results {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(i int, item model.TMDbResult) {
			defer wg.Done()
			defer func() { <-semaphore }()

			detailURL := fmt.Sprintf("%s/%s/%d?api_key=%s", baseURL, mediaType, item.ID, tmdbApiKey)
			// fmt.Println(">>> [DEBUG] Fetching detail:", detailURL)

			detailResp, err := http.Get(detailURL)
			if err != nil {
				fmt.Println(">>> [ERROR] Failed fetch detail:", err)
				return
			}
			defer detailResp.Body.Close()

			detailBody, _ := io.ReadAll(detailResp.Body)

			var detail model.TMDbDetailResponse
			if err := json.Unmarshal(detailBody, &detail); err != nil {
				fmt.Println(">>> [ERROR] Decode detail error:", err)
				return
			}

			movies[i] = MapTMDbToMovieDetailsWithDetail(item, mediaType, detail)
		}(i, item)
	}

	wg.Wait()

	return movies, nil
}

func MapTMDbToMovieDetailsWithDetail(item model.TMDbResult, mediaType string, detail model.TMDbDetailResponse) response.MovieDetailOnlyResponse {
	title := item.Title
	if mediaType == "tv" {
		title = item.Name
	}

	releaseDate := item.ReleaseDate
	if mediaType == "tv" {
		releaseDate = item.FirstAirDate
	}

	genres := make([]string, 0)
	for _, g := range detail.Genres {
		genres = append(genres, g.Name)
	}

	studio := ""
	if len(detail.ProductionCompanies) > 0 {
		studio = detail.ProductionCompanies[0].Name
	}

	totalEps := ""
	if mediaType == "tv" && detail.NumberOfEpisodes > 0 {
		totalEps = strconv.Itoa(detail.NumberOfEpisodes)
	}

	res := response.MovieDetailOnlyResponse{
		MovieID:      "",
		MovieType:    mediaType,
		ThumbnailURL: "https://image.tmdb.org/t/p/w500" + item.PosterPath,
		Title:        title,
		Rating:       fmt.Sprintf("%.1f", item.VoteAverage),
		Status:       detail.Status,
		TotalEps:     totalEps,
		Studio:       studio,
		ReleaseDate:  releaseDate,
		Synopsis:     item.Overview,
		Genres:       genres,
	}

	return res
}

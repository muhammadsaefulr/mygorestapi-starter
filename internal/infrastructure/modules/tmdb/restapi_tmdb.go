package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func FetchTMDbDetail(id int, mediaType string, withRekom bool) (response.MovieDetailOnlyResponse, error) {
	baseURL := "https://api.themoviedb.org/3"
	apiKey := tmdbApiKey

	detailURL := fmt.Sprintf("%s/%s/%d?api_key=%s&language=en-US", baseURL, mediaType, id, apiKey)
	detailResp, err := http.Get(detailURL)
	if err != nil {
		return response.MovieDetailOnlyResponse{}, fmt.Errorf("failed to fetch detail: %w", err)
	}
	defer detailResp.Body.Close()

	var detail model.TMDbDetailResponse
	if err := json.NewDecoder(detailResp.Body).Decode(&detail); err != nil {
		return response.MovieDetailOnlyResponse{}, fmt.Errorf("decode detail failed: %w", err)
	}
	main := MapTMDbToMovieDetailsWithDetail(model.TMDbResult{ID: id}, mediaType, detail)

	rekomURL := fmt.Sprintf("%s/%s/%d/recommendations?api_key=%s&language=en-US&page=1", baseURL, mediaType, id, apiKey)
	rekomResp, err := http.Get(rekomURL)
	if err != nil {
		log.Printf("⚠️ TMDb recommendation fetch error: %v", err)
		return main, nil
	}
	defer rekomResp.Body.Close()

	var rekomResult model.TMDbResponse
	if err := json.NewDecoder(rekomResp.Body).Decode(&rekomResult); err != nil {
		log.Printf("⚠️ TMDb recommendation decode error: %v", err)
		return main, nil
	}

	if len(rekomResult.Results) > 0 {
		rekomID := rekomResult.Results[0].ID

		rekomDetailURL := fmt.Sprintf("%s/%s/%d?api_key=%s", baseURL, mediaType, rekomID, apiKey)
		rekomDetailResp, err := http.Get(rekomDetailURL)
		if err == nil {
			defer rekomDetailResp.Body.Close()
			var rekomDetail model.TMDbDetailResponse
			if err := json.NewDecoder(rekomDetailResp.Body).Decode(&rekomDetail); err == nil {
				rekom := MapTMDbToMovieDetailsWithDetail(rekomResult.Results[0], mediaType, rekomDetail)
				main.Rekomend = &[]response.MovieDetailOnlyResponse{rekom}
			}
		}
	}

	return main, nil
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

	// log.Println(item.ID)

	// strconv.Itoa(item.ID)

	res := response.MovieDetailOnlyResponse{
		IDSource:     strconv.Itoa(item.ID),
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

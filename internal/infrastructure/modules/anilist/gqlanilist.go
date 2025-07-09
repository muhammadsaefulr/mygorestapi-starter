package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/anilist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
)

func FetchAniListMedia(category, queryTitle string, param *request.QueryAnilist) ([]response.MovieDetailOnlyResponse, error) {
	endpoint := "https://graphql.anilist.co"

	variables := map[string]interface{}{
		"page":    param.Page,
		"perPage": param.Limit,
		"type":    "ANIME",
	}

	// log.Printf("Param: %+v", param)

	switch category {
	case "popular":
		variables["sort"] = []string{"POPULARITY_DESC"}
	case "trending":
		variables["sort"] = []string{"TRENDING_DESC"}
	case "ongoing":
		variables["status"] = "RELEASING"
	case "search":
		variables["search"] = param.Search
	default:
		return nil, fmt.Errorf("invalid category")
	}

	query := `query ($page: Int, $perPage: Int, $type: MediaType, $sort: [MediaSort], $search: String, $status: MediaStatus) {
		Page(page: $page, perPage: $perPage) {
			media(type: $type, sort: $sort, search: $search, status: $status) {
				id
				title { romaji english native }
				coverImage { large }
				averageScore
				genres
				status
				episodes
				studios(isMain: true) { nodes { name } }
				startDate { year month day }
				description
			}
		}
	}`

	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result model.AniListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var movies []response.MovieDetailOnlyResponse
	for _, m := range result.Data.Page.Media {
		movies = append(movies, MapAniListToMovieDetails(m, "anime"))
	}

	return movies, nil
}

func FetchAniListDetail(id string) (response.MovieDetailOnlyResponse, error) {
	query := `
	query ($id: Int) {
		Media(id: $id, type: ANIME) {
			id
			title { romaji english native }
			coverImage { large }
			averageScore
			genres
			status
			episodes
			studios(isMain: true) { nodes { name } }
			startDate { year month day }
			description

			recommendations(sort: RATING_DESC, page: 1, perPage: 4) {
				nodes {
					mediaRecommendation {
						id
						title { romaji english native }
						coverImage { large }
						averageScore
						genres
						status
						episodes
						studios(isMain: true) { nodes { name } }
						startDate { year month day }
						description
					}
				}
			}

		}
	}`

	reqBody := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"id": id,
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	resp, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return response.MovieDetailOnlyResponse{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Media struct {
				model.AniListMedia
				Recommendations struct {
					Nodes []struct {
						MediaRecommendation model.AniListMedia `json:"mediaRecommendation"`
					} `json:"nodes"`
				} `json:"recommendations"`
			} `json:"Media"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return response.MovieDetailOnlyResponse{}, err
	}

	main := MapAniListToMovieDetails(result.Data.Media.AniListMedia, "anime")

	if len(result.Data.Media.Recommendations.Nodes) > 0 {
		var rekomendasiList []response.MovieDetailOnlyResponse
		for _, node := range result.Data.Media.Recommendations.Nodes {
			rekom := MapAniListToMovieDetails(node.MediaRecommendation, "anime")
			rekomendasiList = append(rekomendasiList, rekom)
		}
		main.Rekomend = &rekomendasiList
	}

	return main, nil
}

func MapAniListToMovieDetails(m model.AniListMedia, movieType string) response.MovieDetailOnlyResponse {
	releaseDate := fmt.Sprintf("%d-%02d-%02d", m.StartDate.Year, m.StartDate.Month, m.StartDate.Day)
	studio := ""
	if len(m.Studios.Nodes) > 0 {
		studio = m.Studios.Nodes[0].Name
	}
	return response.MovieDetailOnlyResponse{
		IDSource:     strconv.Itoa(m.ID),
		MovieID:      "",
		MovieType:    movieType,
		ThumbnailURL: m.CoverImage.Large,
		Title:        m.Title.Romaji,
		Rating:       strconv.Itoa(m.AverageScore),
		Status:       m.Status,
		TotalEps:     strconv.Itoa(m.Episodes),
		Studio:       studio,
		ReleaseDate:  releaseDate,
		Synopsis:     m.Description,
		Genres:       m.Genres,
	}
}

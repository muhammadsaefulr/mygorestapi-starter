package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/anilist/request"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
)

func FetchAniListMedia(category, queryTitle string, param *request.QueryAnilist) ([]response.MovieDetailOnlyResponse, int, error) {
	endpoint := "https://graphql.anilist.co"

	variables := map[string]interface{}{
		"page":    param.Page,
		"perPage": param.Limit + 3,
		"type":    "ANIME",
	}

	switch category {
	case "popular":
		variables["sort"] = []string{"POPULARITY_DESC"}
	case "trending":
		variables["sort"] = []string{"TRENDING_DESC"}
	case "ongoing":
		season, seasonYear := GetCurrentSeason()
		variables["season"] = season
		variables["seasonYear"] = seasonYear
		variables["status_in"] = []string{"RELEASING", "NOT_YET_RELEASED"}
		variables["format_in"] = []string{"TV", "MOVIE", "ONA", "OVA"}
		variables["sort"] = []string{"POPULARITY_DESC"}
	case "search":
		variables["search"] = param.Search
	case "genre":
		variables["genre_in"] = []string{param.Genre}
		variables["sort"] = []string{"POPULARITY_DESC"}
	default:
		return nil, 0, fmt.Errorf("invalid category")
	}

	query := `query (
		$page: Int, 
		$perPage: Int, 
		$type: MediaType, 
		$sort: [MediaSort], 
		$search: String, 
		$status: MediaStatus,  
		$genre_in: [String],
		$season: MediaSeason,
		$seasonYear: Int,
		$status_in: [MediaStatus],
		$format_in: [MediaFormat]
	) {
		Page(page: $page, perPage: $perPage) {
			pageInfo {
				total
				perPage
				currentPage
				lastPage
				hasNextPage
			}
			media(
				type: $type,
				sort: $sort,
				search: $search,
				status: $status,
				genre_in: $genre_in,
				season: $season,
				seasonYear: $seasonYear,
				status_in: $status_in,
				format_in: $format_in
			) {
				id
				title { romaji english native }
				coverImage { large }
				averageScore
				genres
				status
				episodes
				studios(isMain: true) { nodes { name } }
				startDate { year month day }
				nextAiringEpisode {
					airingAt
					episode
				}
				description
			}
		}
	}`

	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	log.Printf("bodyBytes: %s", bodyBytes)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var result model.AniListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	var movies []response.MovieDetailOnlyResponse
	for _, m := range result.Data.Page.Media {
		movies = append(movies, MapAniListToMovieDetails(m, "anime"))
	}

	return movies, result.Data.Page.PageInfo.LastPage, nil
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
						nextAiringEpisode {
							airingAt
							episode
						}
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

func GetAniListAllGenres() ([]response.GenreDetail, error) {
	baseURL := "https://graphql.anilist.co"
	query := `
		{
			GenreCollection
		}
	`
	reqBody := map[string]interface{}{
		"query": query,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var genreResult model.AniListGenreResponse
	if err := json.NewDecoder(resp.Body).Decode(&genreResult); err != nil {
		return nil, fmt.Errorf("failed to decode genre response: %w", err)
	}

	var genres []response.GenreDetail
	for _, g := range genreResult.Data.GenreCollection {
		genres = append(genres, response.GenreDetail{
			GenreName: g,
			GenreUrl:  fmt.Sprintf("/discovery?type=anime&genre=%s&category=genre", strings.ToLower(strings.ReplaceAll(g, " ", "-"))),
		})
	}

	return genres, nil
}

func MapAniListToMovieDetails(m model.AniListMedia, movieType string) response.MovieDetailOnlyResponse {
	releaseDate := fmt.Sprintf("%d-%02d-%02d", m.StartDate.Year, m.StartDate.Month, m.StartDate.Day)
	studio := ""
	if len(m.Studios.Nodes) > 0 {
		studio = m.Studios.Nodes[0].Name
	}

	log.Printf("AiringDate: %+v", m.NextAiring)

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
		UpdateDay:    utils.GetDayByTimestamp(m.NextAiring.AiringAt),
		Synopsis:     m.Description,
		Genres:       m.Genres,
	}
}

func GetCurrentSeason() (string, int) {
	now := time.Now()
	month := now.Month()
	year := now.Year()

	var season string
	switch {
	case month >= 1 && month <= 3:
		season = "WINTER"
	case month >= 4 && month <= 6:
		season = "SPRING"
	case month >= 7 && month <= 9:
		season = "SUMMER"
	default:
		season = "FALL"
	}
	return season, year
}

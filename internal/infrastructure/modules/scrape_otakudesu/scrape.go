package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func ScrapeHomePage() []model.AnimeData {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0"))
	var results []model.AnimeData

	c.OnHTML(".venz li", func(e *colly.HTMLElement) {
		results = append(results, model.AnimeData{
			Title:        e.ChildText(".jdlflm"),
			URL:          e.ChildAttr(".thumb a", "href"),
			JudulPath:    strings.TrimSuffix(path.Base(e.ChildAttr(".thumb a", "href")), "/"),
			ThumbnailURL: e.ChildAttr(".thumb img", "src"),
			LatestEp:     e.ChildText(".epz"),
			UpdateAnime:  e.ChildText(".epztipe"),
		})
	})
	_ = c.Visit("https://otakudesu.cloud/")
	return results
}

func ScrapeCompleteAnime(page string) []model.CompleteAnime {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0"))
	var results []model.CompleteAnime

	c.OnHTML(".venz li", func(e *colly.HTMLElement) {
		animeURL := e.ChildAttr(".thumb a", "href")
		ratingStr := e.ChildText(".epztipe")

		results = append(results, model.CompleteAnime{
			Title:        e.ChildText(".jdlflm"),
			URL:          animeURL,
			JudulPath:    strings.TrimSuffix(path.Base(animeURL), "/"),
			ThumbnailURL: e.ChildAttr(".thumb img", "src"),
			LatestEp:     e.ChildText(".epz"),
			Rating:       ratingStr,
			UpdatedAt:    e.ChildText(".newnime"),
		})
	})

	_ = c.Visit("https://otakudesu.cloud/complete-anime/page/" + page)

	sort.SliceStable(results, func(i, j int) bool {
		ratingI, _ := strconv.ParseFloat(results[i].Rating, 64)
		ratingJ, _ := strconv.ParseFloat(results[j].Rating, 64)
		return ratingI > ratingJ
	})

	return results
}

func ScrapeOngoingAnime(page string) []model.OngoingAnime {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0"))
	var results []model.OngoingAnime

	c.OnHTML(".venz li", func(e *colly.HTMLElement) {
		animeURL := e.ChildAttr(".thumb a", "href")

		results = append(results, model.OngoingAnime{
			Title:        e.ChildText(".jdlflm"),
			URL:          animeURL,
			JudulPath:    strings.TrimSuffix(path.Base(animeURL), "/"),
			ThumbnailURL: e.ChildAttr(".thumb img", "src"),
			Episode:      e.ChildText(".epz"),
			DaysUpdated:  e.ChildText(".epztipe"),
			UpdatedAt:    e.ChildText(".newnime"),
		})
	})

	err := c.Visit("https://otakudesu.cloud/ongoing-anime/page/" + page)
	if err != nil {
		log.Println("Visit error:", err)
	}

	return results
}

func ScrapeGenreAnime(url string) []model.GenreAnime {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0"))
	var results []model.GenreAnime

	c.OnHTML(".col-anime", func(e *colly.HTMLElement) {
		results = append(results, model.GenreAnime{
			Title:        e.ChildText(".col-anime-title a"),
			URL:          "/detail/" + path.Base(strings.TrimSuffix(e.ChildAttr(".col-anime-title a", "href"), "/")),
			Studio:       e.ChildText(".col-anime-studio"),
			ThumbnailURL: e.ChildAttr(".col-anime-cover img", "src"),
			Episodes:     e.ChildText(".col-anime-eps"),
			Rating:       e.ChildText(".col-anime-rating"),
		})
	})
	_ = c.Visit(url)
	if len(results) > 20 {
		return results[:20]
	}
	return results
}

func ScrapeGenreList(url string) []model.GenreList {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0"))
	var results []model.GenreList

	c.OnHTML("ul.genres a", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.Text)
		href := strings.TrimSpace(e.Attr("href"))

		if name != "" && href != "" {
			results = append(results, model.GenreList{
				Title: name,
				URL:   "/genre/" + path.Base(strings.TrimSuffix(e.Request.AbsoluteURL(href), "/")) + "/page/1",
			})
		}
	})

	_ = c.Visit(url)
	return results
}

func ScrapeAnimeDetail(url string) (model.AnimeDetail, []model.AnimeEpisode, []model.SearchResult) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
		colly.Async(true),
	)
	var (
		detail   model.AnimeDetail
		episodes []model.AnimeEpisode
	)

	c.OnHTML(".fotoanime", func(e *colly.HTMLElement) {
		infozingle := e.DOM.Find(".infozingle")
		sinopc := e.DOM.Find(".sinopc")

		var genres []model.GenreInfo
		infozingle.Find("span b").Each(func(_ int, el *goquery.Selection) {
			if strings.Contains(el.Text(), "Genre") {
				el.Parent().Find("a").Each(func(_ int, s *goquery.Selection) {

					href := s.AttrOr("href", "")
					slug := path.Base(strings.TrimSuffix(href, "/"))
					genres = append(genres, model.GenreInfo{
						Title: s.Text(),
						URL:   "/discovery/genres/" + slug,
					})
				})
			}
		})

		re := regexp.MustCompile(`(?i)\(Info: Episode sebelumnya akan ditambahkan secara berkala\)`)
		cleanSynopsis := re.ReplaceAllString(sinopc.Find("p").Text(), "")

		detail = model.AnimeDetail{
			ThumbnailURL: e.DOM.Find("img").First().AttrOr("src", ""),
			Title:        strings.TrimPrefix(infozingle.Find("p:contains('Judul')").Text(), "Judul: "),
			Rating:       strings.TrimPrefix(infozingle.Find("p:contains('Skor')").Text(), "Skor: "),
			Producer:     strings.TrimPrefix(infozingle.Find("p:contains('Produser')").Text(), "Produser: "),
			Status:       strings.TrimPrefix(infozingle.Find("p:contains('Status')").Text(), "Status: "),
			TotalEps:     strings.TrimPrefix(infozingle.Find("p:contains('Total Episode')").Text(), "Total Episode: "),
			Duration:     strings.TrimSuffix(strings.TrimPrefix(infozingle.Find("p:contains('Durasi')").Text(), "Durasi: "), "per ep."),
			Studio:       strings.TrimPrefix(infozingle.Find("p:contains('Studio:')").Text(), "Studio: "),
			ReleaseDate:  strings.TrimSpace(strings.TrimPrefix(infozingle.Find("p:contains('Tanggal Rilis')").Text(), "Tanggal Rilis:")),
			Synopsis:     cleanSynopsis,
			Genres:       genres,
		}
	})

	c.OnHTML(".episodelist li", func(e *colly.HTMLElement) {
		title := e.ChildText("span a")
		href := e.ChildAttr("span a", "href")
		releaseDate := e.ChildText("span.zeebr")

		if !strings.Contains(strings.ToLower(href), "batch") && !strings.Contains(strings.ToLower(href), "lengkap") {
			episodes = append(episodes, model.AnimeEpisode{
				Title:       title,
				VideoURL:    "/otakudesu/play/" + path.Base(strings.TrimSuffix(href, "/")),
				ReleaseDate: releaseDate,
			})
		}
	})

	_ = c.Visit(url)
	c.Wait()

	recommendations, _ := fetchRecommendationsFromAniList(detail.Title)

	if strings.ToLower(detail.Status) == "ongoing" {
		if len(episodes) >= 2 {
			day1 := utils.ConvertDateStrToDay(episodes[0].ReleaseDate)
			day2 := utils.ConvertDateStrToDay(episodes[1].ReleaseDate)
			log.Printf("Day1: %s, Day2: %s", episodes[len(episodes)-2].ReleaseDate, episodes[len(episodes)-1].ReleaseDate)

			if day1 == day2 {
				detail.UpdatedDay = day2
			} else {
				detail.UpdatedDay = day2
			}
		} else if len(episodes) == 1 {
			detail.UpdatedDay = utils.ConvertDateStrToDay(episodes[0].ReleaseDate)
		}
	}

	return detail, episodes, recommendations
}

func ScrapeSearchAnimeByTitle(url string) []model.SearchResult {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0"))
	var results []model.SearchResult

	// log.Printf("Scraping search title: %s", url)

	// searchTitle := strings.ReplaceAll(title, " ", "+")

	c.OnHTML("ul.chivsrc li", func(e *colly.HTMLElement) {
		var genres []model.GenreInfo
		e.ForEach(".set b", func(_ int, el *colly.HTMLElement) {
			if strings.Contains(el.Text, "Genres") {
				el.DOM.Parent().Find("a").Each(func(_ int, s *goquery.Selection) {
					genres = append(genres, model.GenreInfo{
						Title: s.Text(),
						URL:   "/discovery/genres/" + path.Base(strings.TrimSuffix(s.AttrOr("href", ""), "/"))})
				})
			}
		})
		results = append(results, model.SearchResult{
			Title:        e.ChildText("h2 a"),
			URL:          "/otakudesu/detail/" + path.Base(strings.TrimSuffix(e.ChildAttr("h2 a", "href"), "/")),
			ThumbnailURL: e.ChildAttr("img", "src"),
			Genres:       genres,
			Status:       strings.TrimSpace(strings.Split(e.DOM.Find(".set b:contains('Status')").Parent().Text(), ":")[1]),
			Rating:       strings.TrimSpace(strings.Split(e.DOM.Find(".set b:contains('Rating')").Parent().Text(), ":")[1]),
		})
	})
	_ = c.Visit(url + "&post_type=anime")
	if len(results) > 15 {
		return results[:15]
	}
	return results
}

func ScrapeAnimeSourceData(url string) model.AnimeSourceData {
	c := colly.NewCollector()
	var epsList []model.AnimeEpisode
	var animeSource []model.VideoSource
	var result model.AnimeSourceData

	c.OnHTML(".keyingpost li", func(e *colly.HTMLElement) {
		epsList = append(epsList, model.AnimeEpisode{
			Title:    e.ChildText("a"),
			VideoURL: "/otakudesu/play/" + path.Base(strings.TrimSuffix(e.ChildAttr("a", "href"), "/")),
		})
	})

	c.OnHTML(".download ul li", func(e *colly.HTMLElement) {
		var dataList []model.AnimeEpisode
		titleRes := e.ChildText("strong")

		var wg sync.WaitGroup
		var mu sync.Mutex

		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			title := strings.TrimSpace(el.Text)
			link := el.Attr("href")

			if strings.EqualFold(title, "pdrain") {
				wg.Add(1)
				go func(title, link string) {
					defer wg.Done()
					extracted := ExtractPdrainUrl(link)
					if extracted != "" {
						mu.Lock()
						dataList = append(dataList, model.AnimeEpisode{
							Title:    title,
							VideoURL: extracted,
						})
						mu.Unlock()
					}
				}(title, link)
			}
			//  else {
			// 	mu.Lock()
			// 	dataList = append(dataList, model.AnimeEpisode{
			// 		Title:    title,
			// 		VideoURL: link,
			// 	})
			// 	mu.Unlock()
			// }
		})

		wg.Wait()

		animeSource = append(animeSource, model.VideoSource{
			Res:      titleRes,
			DataList: dataList,
		})
	})

	c.OnHTML(".venutama h1.posttl", func(e *colly.HTMLElement) {
		txt := e.Text

		epRegex := regexp.MustCompile(`(?i)(Episode\s\d+)`)
		if match := epRegex.FindString(txt); match != "" {
			result.CurrentEp = strings.TrimSpace(match)
		}

		titleRegex := regexp.MustCompile(`(?i)\s+Episode\s\d+.*`)
		result.Title = strings.TrimSpace(titleRegex.ReplaceAllString(txt, ""))
	})

	c.OnHTML(".kategoz span", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Release on") {
			result.ReleaseDate = strings.TrimSpace(strings.Replace(e.Text, "Release on", "", 1))
		}
	})

	c.OnHTML(".cukder img", func(e *colly.HTMLElement) {
		result.ThumbnailURL = e.Attr("src")
	})

	c.OnHTML(".responsive-embed-stream iframe", func(e *colly.HTMLElement) {
		result.DownloadURL = e.Attr("src")
	})

	c.OnHTML(".flir a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Next Eps.") {
			result.NextEpURL = e.Attr("href")
		}

		if strings.Contains(e.Text, "See All Episodes") {
			result.DetailURL = "/detail/" + path.Base(strings.TrimSuffix(e.Attr("href"), "/"))
		}
	})

	if err := c.Visit(url); err != nil {
		log.Println("Visit error:", err)
	}

	result.Episodes = epsList
	result.Sources = animeSource
	return result
}

// Utils

func ExtractPdrainUrl(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching Pdrain URL:", err)
		return ""
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("Error parsing Pdrain HTML:", err)
		return ""
	}

	return doc.Find(`meta[name="twitter:player:stream"]`).AttrOr("content", "")
}

func fetchRecommendationsFromAniList(animeURL string) ([]model.SearchResult, error) {

	query := `
	query ($search: String) {
		Media(search: $search, type: ANIME) {
			recommendations(page: 1, perPage: 10, sort: RATING_DESC) {
				nodes {
					mediaRecommendation {
						title {
							romaji
						}
						siteUrl
						genres
						averageScore
						coverImage {
							large
						}
					}
				}
			}
		}
	}`

	variables := map[string]interface{}{
		"search": animeURL,
	}
	body := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://graphql.anilist.co", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	type gqlResponse struct {
		Data struct {
			Media struct {
				Recommendations struct {
					Nodes []struct {
						MediaRecommendation struct {
							Title struct {
								Romaji string `json:"romaji"`
							} `json:"title"`
							SiteURL      string   `json:"siteUrl"`
							Genres       []string `json:"genres"`
							AverageScore int      `json:"averageScore"`
							CoverImage   struct {
								Large string `json:"large"`
							} `json:"coverImage"`
						} `json:"mediaRecommendation"`
					} `json:"nodes"`
				} `json:"recommendations"`
			} `json:"Media"`
		} `json:"data"`
	}

	var gqlRes gqlResponse
	if err := json.Unmarshal(respBody, &gqlRes); err != nil {
		return nil, err
	}

	var (
		results []model.SearchResult
		mu      sync.Mutex
		wg      sync.WaitGroup
		sem     = make(chan struct{}, 6) // max pool
	)

	for _, node := range gqlRes.Data.Media.Recommendations.Nodes {
		m := node.MediaRecommendation

		wg.Add(1)
		go func(mrTitle string, mGenres []string, mCover string, mScore int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			searchResults := ScrapeSearchAnimeByTitle("https://otakudesu.cloud/?s=" + url.QueryEscape(mrTitle))
			if len(searchResults) == 0 {
				// log.Println("[INFO] Gagal temukan di Otakudesu:", mrTitle)
				return
			}

			first := searchResults[0]

			if first.ThumbnailURL == "" {
				first.ThumbnailURL = mCover
			}
			if first.Rating == "" && mScore > 0 {
				first.Rating = fmt.Sprintf("%.2f", float64(mScore)/10)
			}
			if len(first.Genres) == 0 && len(mGenres) > 0 {
				for _, g := range mGenres {
					first.Genres = append(first.Genres, model.GenreInfo{
						Title: strings.ToLower(g),
					})
				}
			}
			if first.Status == "" {
				first.Status = "Unknown"
			}

			mu.Lock()
			results = append(results, first)
			mu.Unlock()
		}(m.Title.Romaji, m.Genres, m.CoverImage.Large, m.AverageScore)
	}

	wg.Wait()
	return results, nil
}

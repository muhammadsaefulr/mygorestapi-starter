package modules

import (
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"

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

func ScrapeGenreAnime(url string) []model.GenreAnime {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0"))
	var results []model.GenreAnime

	c.OnHTML(".col-anime", func(e *colly.HTMLElement) {
		results = append(results, model.GenreAnime{
			Title:    e.ChildText(".col-anime-title a"),
			URL:      e.ChildAttr(".col-anime-title a", "href"),
			Studio:   e.ChildText(".col-anime-studio"),
			Episodes: e.ChildText(".col-anime-eps"),
			Rating:   e.ChildText(".col-anime-rating"),
		})
	})
	_ = c.Visit(url)
	if len(results) > 20 {
		return results[:20]
	}
	return results
}

func ScrapeAnimeEpisodes(url string) (model.AnimeDetail, []model.AnimeEpisode) {
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
					genres = append(genres, model.GenreInfo{
						Title: s.Text(),
						URL:   s.AttrOr("href", ""),
					})
				})
			}
		})

		re := regexp.MustCompile(`(?i)\(Info: Episode sebelumnya akan ditambahkan secara berkala\)`)
		cleanSynopsis := re.ReplaceAllString(sinopc.Find("p").Text(), "")

		detail = model.AnimeDetail{
			ThumbnailURL: infozingle.Parent().Find("img.attachment-post-thumbnail.size-post-thumbnail").AttrOr("src", ""),
			Title:        strings.TrimPrefix(infozingle.Find("p:contains('Judul')").Text(), "Judul: "),
			Rating:       strings.TrimPrefix(infozingle.Find("p:contains('Skor')").Text(), "Skor: "),
			Producer:     strings.TrimPrefix(infozingle.Find("p:contains('Produser')").Text(), "Produser: "),
			Status:       strings.TrimPrefix(infozingle.Find("p:contains('Status')").Text(), "Status: "),
			TotalEps:     strings.TrimPrefix(infozingle.Find("p:contains('Total Episode')").Text(), "Total Episode: "),
			Duration:     strings.TrimSuffix(strings.TrimPrefix(infozingle.Find("p:contains('Durasi')").Text(), "Durasi: "), "per ep."),
			Studio:       strings.TrimPrefix(infozingle.Find("p:contains('Studio:')").Text(), "Studio: "),
			ReleaseDate:  infozingle.Find("p:contains('Tanggal Rilis')").Text(),
			Synopsis:     cleanSynopsis,
			Genres:       genres,
		}
	})

	c.OnHTML(".episodelist li", func(e *colly.HTMLElement) {
		episodes = append(episodes, model.AnimeEpisode{
			Title:    e.ChildText("span a"),
			VideoURL: e.ChildAttr("span a", "href"),
		})
	})

	_ = c.Visit(url)
	c.Wait()

	return detail, episodes
}

func ScrapeSearchAnimeByTitle(url string) []model.SearchResult {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0"))
	var results []model.SearchResult

	c.OnHTML("ul.chivsrc li", func(e *colly.HTMLElement) {
		var genres []model.GenreInfo
		e.ForEach(".set b", func(_ int, el *colly.HTMLElement) {
			if strings.Contains(el.Text, "Genres") {
				el.DOM.Parent().Find("a").Each(func(_ int, s *goquery.Selection) {
					genres = append(genres, model.GenreInfo{
						Title: s.Text(),
						URL:   s.AttrOr("href", ""),
					})
				})
			}
		})
		results = append(results, model.SearchResult{
			Title:        e.ChildText("h2 a"),
			URL:          e.ChildAttr("h2 a", "href"),
			ThumbnailURL: e.ChildAttr("img", "src"),
			Genres:       genres,
			Status:       e.ChildText(".set b:contains('Status')"),
			Rating:       e.ChildText(".set b:contains('Rating')"),
		})
	})
	_ = c.Visit(url)
	if len(results) > 15 {
		return results[:15]
	}
	return results
}

func ScrapeOngoingAnime(url string) []model.AnimeData {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0"))
	var results []model.AnimeData

	c.OnHTML(".venz li", func(e *colly.HTMLElement) {
		results = append(results, model.AnimeData{
			Title:        e.ChildText(".jdlflm"),
			LatestEp:     e.ChildText(".epz"),
			URL:          e.ChildAttr(".thumb a", "href"),
			UpdateAnime:  e.ChildText(".epztipe"),
			ThumbnailURL: e.ChildAttr(".thumbz img", "src"),
		})
	})
	_ = c.Visit(url)
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
			VideoURL: e.ChildAttr("a", "href"),
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
			} else {
				mu.Lock()
				dataList = append(dataList, model.AnimeEpisode{
					Title:    title,
					VideoURL: link,
				})
				mu.Unlock()
			}
		})

		wg.Wait()

		animeSource = append(animeSource, model.VideoSource{
			Res:      titleRes,
			DataList: dataList,
		})
	})

	c.OnHTML(".venutama h1.posttl", func(e *colly.HTMLElement) {
		txt := e.Text
		r := regexp.MustCompile(`(?i)^(.*?)\s(Episode\s\d+\sSubtitle\sIndonesia)$`)
		if parts := r.FindStringSubmatch(txt); len(parts) >= 3 {
			result.Title = parts[1]
			result.CurrentEp = parts[2]
		}
	})

	c.OnHTML(".kategoz span", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Release on") {
			result.ReleaseDate = strings.TrimSpace(strings.Replace(e.Text, "Release on", "", 1))
		}
	})

	c.OnHTML(".responsive-embed-stream iframe", func(e *colly.HTMLElement) {
		result.DownloadURL = e.Attr("src")
	})

	c.OnHTML(".flir a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Next Eps.") {
			result.NextEpURL = e.Attr("href")
		}
	})

	if err := c.Visit(url); err != nil {
		log.Println("Visit error:", err)
	}

	result.Episodes = epsList
	result.Sources = animeSource
	return result
}

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

// func main() {
// // Test scrapeHomePage
// home := ScrapeHomePage()
// fmt.Println("\n== Home Page ==")
// for _, a := range home {
// 	fmt.Printf("%s | %s\n", a.Title, a.URL)
// }

// // Test scrapemodel.GenreAnimes
// model.GenreAnimes := Scrapemodel.GenreAnimes("https://otakudesu.cloud/genre/romance/")
// fmt.Println("\n== Genre Animes ==")
// for _, g := range model.GenreAnimes {
// 	fmt.Printf("%s | %s\n", g.Title, g.Links)
// }

// Test scrapemodel.AnimeEpisodes
// detail, episodes := ScrapeAnimeEpisodes("https://otakudesu.cloud/anime/zatsu-tabi-journey-sub-indo")
// fmt.Println("\n== Anime Detail ==")
// fmt.Printf("Title: %s\nEpisodes: %d\n", detail.Title, len(episodes))
// for _, e := range episodes {
// 	fmt.Printf("Ep: %s | %s\n", e.Title, e.VideoURL)
// }

// sourceData := scrapeAnimeSourceData("https://otakudesu.cloud/episode/zttj-episode-2-sub-indo/")

// fmt.Println("\n== Anime Source Data ==")
// fmt.Println(sourceData)

// Test scrapeSearchAnimeByTitle
// model.SearchResults := ScrapeSearchAnimeByTitle("https://otakudesu.cloud/?s=one+piece")
// fmt.Println("\n== Search Result ==")
// for _, s := range model.SearchResults {
// 	fmt.Printf("%s | %s\n", s.Title, s.AnimeLinks)
// }

// // Test scrapeOngoingAnime
// ongoing := ScrapeOngoingAnime("https://otakudesu.cloud/ongoing-anime/")
// fmt.Println("\n== Ongoing Anime ==")
// for _, o := range ongoing {
// 	fmt.Printf("%s | %s\n", o.Title, o.URL)
// }

// 	fmt.Println("\n== DONE TESTING ==")
// }

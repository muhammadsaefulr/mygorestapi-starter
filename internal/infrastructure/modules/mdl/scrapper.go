package scrapper

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	// "sync"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/dto/movie_details/response"
)

func NewChromeContext() (context.Context, context.CancelFunc) {
	allocatorCtx, cancelAllocator := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
				"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"),
		)...,
	)
	ctx, cancelCtx := chromedp.NewContext(allocatorCtx)
	ctx, timeoutCancel := context.WithTimeout(ctx, 60*time.Second)

	return ctx, func() {
		timeoutCancel()
		cancelCtx()
		cancelAllocator()
	}
}

func GetDramaDetail(parentCtx context.Context, urlstr string) (response.MovieDetailOnlyResponse, error) {
	ctx, cancel := chromedp.NewContext(parentCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		network.Enable(),
		network.SetBlockedURLs([]string{"*.png", "*.jpg", "*.jpeg", "*.css", "*.js", "*.woff"}),
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible("#show-detailsxx", chromedp.ByID),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return response.MovieDetailOnlyResponse{}, fmt.Errorf("navigate detail failed: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return response.MovieDetailOnlyResponse{}, fmt.Errorf("parse detail failed: %w", err)
	}

	now := time.Now()

	title := strings.TrimSpace(doc.Find(".film-title-wrapper").Text())
	if title == "" {
		title = strings.TrimSpace(doc.Find("h1.film-title").Text())
	}

	thumbnail, _ := doc.Find(".film-cover img").Attr("src")
	rating := strings.TrimSpace(doc.Find(`.hfs > b[itempropx="ratingValue"]`).Text())
	synopsis := strings.TrimSpace(doc.Find(".show-synopsis > p").Text())

	status := "Unknown"
	totalEps := "Unknown"
	aired := "Unknown"
	studio := "Unknown"
	producer := "Unknown"
	var genres []string

	doc.Find(".show-detailsxss ul.list li").Each(func(i int, s *goquery.Selection) {
		key := strings.TrimSpace(s.Find("b.inline").Text())

		switch key {
		case "Status:":
			status = strings.TrimSpace(s.Text()[len(key):])
		case "Episodes:":
			totalEps = strings.TrimSpace(s.Text()[len(key):])
		case "Director:":
			producer = strings.TrimSpace(s.Text()[len(key):])
		case "Aired:":
			text := strings.TrimSpace(s.Text()[len(key):])
			dates := strings.Split(text, "-")
			if len(dates) > 0 {
				start := strings.TrimSpace(dates[0])
				aired = start
			}
			if len(dates) == 2 {
				endStr := strings.TrimSpace(dates[1])
				layout := "Jan 2, 2006"

				if t, err := time.Parse(layout, endStr); err == nil {
					if t.Before(time.Now()) {
						status = "Completed"
					} else {
						status = "Ongoing"
					}
				}
			}
		case "Original Network:":
			studio = strings.TrimSpace(s.Find("a").First().Text())
		case "Genres:":
			s.Find("a").Each(func(i int, g *goquery.Selection) {
				genres = append(genres, strings.TrimSpace(g.Text()))
			})
		}
	})

	return response.MovieDetailOnlyResponse{
		MovieType:    "drama",
		Title:        title,
		Rating:       rating,
		Synopsis:     synopsis,
		Status:       status,
		TotalEps:     totalEps,
		ReleaseDate:  aired,
		Studio:       studio,
		Producer:     producer,
		Genres:       genres,
		ThumbnailURL: thumbnail,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}, nil
}

func FetchMDLMedia(ctx context.Context, category string, page, limit int) ([]response.MovieDetailOnlyResponse, int, int, error) {
	var html string

	base := "https://mydramalist.com/search?adv=titles&ty=68&co=3"
	switch category {
	case "trending":
		base += "&so=top"
	case "popular":
		base += "&so=popular"
	case "ongoing":
		base += "&st=1"
	default:
		return nil, 0, 0, fmt.Errorf("kategori tidak dikenali: %s", category)
	}

	urlStr := fmt.Sprintf("%s&page=%d", base, page)

	if err := chromedp.Run(ctx,
		network.Enable(),
		network.SetBlockedURLs([]string{"*.png", "*.jpg", "*.jpeg", "*.css", "*.js", "*.woff"}),
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(".box", chromedp.ByQuery),
		chromedp.OuterHTML("body", &html),
	); err != nil {
		return nil, 0, 0, fmt.Errorf("navigate error: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("parse HTML: %w", err)
	}

	totalPages := 1
	doc.Find("ul.pagination li.page-item.last a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			u, err := url.Parse(href)
			if err == nil {
				q := u.Query()
				if pg := q.Get("page"); pg != "" {
					if n, err := strconv.Atoi(pg); err == nil {
						totalPages = n
					}
				}
			}
		}
	})

	type linkInfo struct {
		MovieID string
		URL     string
	}
	var links []linkInfo

	doc.Find(".box h6.title a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			fullURL := "https://mydramalist.com" + href
			slug := strings.TrimPrefix(href, "/")
			movieID := strings.Split(slug, "-")[0]
			links = append(links, linkInfo{MovieID: movieID, URL: fullURL})
		}
	})

	if len(links) > limit {
		links = links[:limit]
	}

	const maxGoroutines = 6
	var (
		results []response.MovieDetailOnlyResponse
		wg      sync.WaitGroup
		mutex   sync.Mutex
		sem     = make(chan struct{}, maxGoroutines)
	)

	for _, link := range links {
		wg.Add(1)
		go func(link linkInfo) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			detail, err := GetDramaDetail(ctx, link.URL)
			if err != nil {
				log.Printf("❌ gagal ambil detail %s: %v", link.URL, err)
				return
			}

			detail.MovieID = link.MovieID
			detail.MovieType = "drama"

			mutex.Lock()
			results = append(results, detail)
			mutex.Unlock()
		}(link)
	}

	wg.Wait()
	return results, len(results), totalPages, nil
}

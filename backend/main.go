package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// Structure pour stocker les informations sur chaque lien
type LinkInfo struct {
	URL                  string `json:"url"`
	Title                string `json:"title"`
	Status               int    `json:"status"`
	Visited              bool   `json:"visited"`
	HasH1                bool   `json:"hasH1"`
	MultipleH1s          bool   `json:"multipleH1s"`
	MissingAlts          int    `json:"missingAlts"`
	TotalImages          int    `json:"totalImages"`
	ButtonsWithoutLabels int    `json:"buttonsWithoutLabels"`
	HasMetaDescription   bool   `json:"hasMetaDescription"`
	MetaDescription      string `json:"metaDescription"`
	HasCanonical         bool   `json:"hasCanonical"`
	HasRobotsMeta        bool   `json:"hasRobotsMeta"`
	RobotsMetaValue      string `json:"robotsMetaValue"`
	HasNoindex           bool   `json:"hasNoindex"`
	HasFavicon           bool   `json:"hasFavicon"`
	HasOpenGraph         bool   `json:"hasOpenGraph"`
	HasTwitterCard       bool   `json:"hasTwitterCard"`
	HasStructuredData    bool   `json:"hasStructuredData"`
	HasViewport          bool   `json:"hasViewport"`
	HasHtmlLang          bool   `json:"hasHtmlLang"`
	TitleEmptyOrShort    bool   `json:"titleEmptyOrShort"`
	MetaDescEmptyOrShort bool   `json:"metaDescEmptyOrShort"`
}

type ScrapeResult struct {
	TotalLinks int                 `json:"totalLinks"`
	Links      []*LinkInfo         `json:"links"`
	Stats      map[string]int      `json:"stats"`
}

// Map pour stocker tous les liens trouvés
var allLinks = make(map[string]*LinkInfo)
var mutex sync.RWMutex

func analyzeSEO(url string) (LinkInfo, error) {
	var info LinkInfo
	info.URL = url
	resp, err := http.Get(url)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return info, err
	}

	// Title
	title := strings.TrimSpace(doc.Find("title").Text())
	info.Title = title
	info.TitleEmptyOrShort = len(title) < 10

	// Meta Description
	metaDesc, exists := doc.Find("meta[name='description']").Attr("content")
	info.HasMetaDescription = exists
	info.MetaDescription = strings.TrimSpace(metaDesc)
	info.MetaDescEmptyOrShort = len(info.MetaDescription) < 30

	// Canonical
	_, info.HasCanonical = doc.Find("link[rel='canonical']").Attr("href")

	// Robots Meta
	robotsMeta, robotsExists := doc.Find("meta[name='robots']").Attr("content")
	info.HasRobotsMeta = robotsExists
	info.RobotsMetaValue = strings.TrimSpace(robotsMeta)
	info.HasNoindex = strings.Contains(strings.ToLower(info.RobotsMetaValue), "noindex")

	// H1s
	h1Count := doc.Find("h1").Length()
	info.HasH1 = h1Count > 0
	info.MultipleH1s = h1Count > 1

	// Images and Alts
	info.TotalImages = doc.Find("img").Length()
	info.MissingAlts = 0
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		if _, exists := s.Attr("alt"); !exists {
			info.MissingAlts++
		}
	})

	// Buttons without labels
	info.ButtonsWithoutLabels = 0
	doc.Find("button").Each(func(i int, s *goquery.Selection) {
		title, tExists := s.Attr("title")
		aria, aExists := s.Attr("aria-label")
		if (!tExists || strings.TrimSpace(title) == "") && (!aExists || strings.TrimSpace(aria) == "") {
			info.ButtonsWithoutLabels++
		}
	})

	// Favicon
	info.HasFavicon = doc.Find("link[rel='icon'], link[rel='shortcut icon']").Length() > 0

	// Open Graph
	info.HasOpenGraph = doc.Find("meta[property^='og:' ]").Length() > 0

	// Twitter Card
	info.HasTwitterCard = doc.Find("meta[name^='twitter:' ]").Length() > 0

	// Structured Data
	info.HasStructuredData = doc.Find("script[type='application/ld+json']").Length() > 0

	// Viewport
	info.HasViewport = doc.Find("meta[name='viewport']").Length() > 0

	// HTML lang
	lang, exists := doc.Find("html").Attr("lang")
	info.HasHtmlLang = exists && strings.TrimSpace(lang) != ""

	return info, nil
}

func scrapeSite(startURL string, maxDepth int) (*ScrapeResult, error) {
	allLinks := make(map[string]*LinkInfo)
	var mutex sync.RWMutex

	parsedURL, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("URL invalide: %v", err)
	}
	baseDomain := parsedURL.Host

	c := colly.NewCollector(
		colly.AllowedDomains(baseDomain),
		colly.MaxDepth(maxDepth),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
		RandomDelay: 1,
	})

	c.OnRequest(func(r *colly.Request) {
		mutex.Lock()
		if _, exists := allLinks[r.URL.String()]; !exists {
			allLinks[r.URL.String()] = &LinkInfo{
				URL:     r.URL.String(),
				Visited: false,
			}
		}
		mutex.Unlock()
	})

	c.OnResponse(func(r *colly.Response) {
		mutex.Lock()
		if link, exists := allLinks[r.Request.URL.String()]; exists {
			link.Status = r.StatusCode
			link.Visited = true
			if r.StatusCode == 200 {
				seoInfo, err := analyzeSEO(r.Request.URL.String())
				if err == nil {
					*link = seoInfo
					link.Status = r.StatusCode
					link.Visited = true
				}
			}
		}
		mutex.Unlock()
	})

	c.OnError(func(r *colly.Response, err error) {
		mutex.Lock()
		if link, exists := allLinks[r.Request.URL.String()]; exists {
			link.Status = r.StatusCode
			link.Visited = true
		}
		mutex.Unlock()
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		title := strings.TrimSpace(e.Text)
		absoluteURL := e.Request.AbsoluteURL(href)
		if absoluteURL == "" || strings.HasPrefix(href, "javascript:") ||
			strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") ||
			strings.HasPrefix(href, "#") {
			return
		}
		parsedLink, err := url.Parse(absoluteURL)
		if err != nil || parsedLink.Host != baseDomain {
			return
		}
		mutex.Lock()
		if _, exists := allLinks[absoluteURL]; !exists {
			allLinks[absoluteURL] = &LinkInfo{
				URL:   absoluteURL,
				Title: title,
			}
		}
		mutex.Unlock()
		if !allLinks[absoluteURL].Visited {
			c.Visit(absoluteURL)
		}
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		mutex.Lock()
		if link, exists := allLinks[e.Request.URL.String()]; exists {
			link.Title = strings.TrimSpace(e.Text)
		}
		mutex.Unlock()
	})

	err = c.Visit(startURL)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la visite initiale: %v", err)
	}
	c.Wait()

	visited := 0
	success := 0
	errors := 0
	for _, info := range allLinks {
		if info.Visited {
			visited++
			if info.Status == 200 {
				success++
			} else {
				errors++
			}
		}
	}
	stats := map[string]int{
		"visited":    visited,
		"success":    success,
		"errors":     errors,
		"notVisited": len(allLinks) - visited,
	}
	links := make([]*LinkInfo, 0, len(allLinks))
	for _, l := range allLinks {
		links = append(links, l)
	}
	return &ScrapeResult{
		TotalLinks: len(allLinks),
		Links:      links,
		Stats:      stats,
	}, nil
}

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Query().Get("url")
	if urlParam == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}
	result, err := scrapeSite(urlParam, 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/api/scrape", scrapeHandler)
	fmt.Println("SEO Scraper API running on http://localhost:8080/api/scrape?url=...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

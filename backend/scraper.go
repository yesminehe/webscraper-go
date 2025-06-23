package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

var latestScrapeResult *ScrapeResult

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
		return nil, fmt.Errorf("erreur lors de la visite initiale: %v", err)
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
	result := &ScrapeResult{
		TotalLinks: len(allLinks),
		Links:      links,
		Stats:      stats,
	}
	latestScrapeResult = result
	return result, nil
}

// ExportScrapeResultToJSON exports the scrape result to a JSON file.
func ExportScrapeResultToJSON(result *ScrapeResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// ExportScrapeResultToCSV exports the scrape result to a CSV file.
func ExportScrapeResultToCSV(result *ScrapeResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"URL", "Title", "Status", "Visited"})

	for _, link := range result.Links {
		writer.Write([]string{
			link.URL,
			link.Title,
			fmt.Sprintf("%d", link.Status),
			fmt.Sprintf("%v", link.Visited),
		})
	}
	return nil
}

func exportJSONHandler(w http.ResponseWriter, r *http.Request) {
	if latestScrapeResult == nil {
		http.Error(w, "No scrape result available", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=\"scrape_result.json\"")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(latestScrapeResult)
}

func exportCSVHandler(w http.ResponseWriter, r *http.Request) {
	if latestScrapeResult == nil {
		http.Error(w, "No scrape result available", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=\"scrape_result.csv\"")
	w.Header().Set("Content-Type", "text/csv")
	writer := csv.NewWriter(w)
	writer.Write([]string{"URL", "Title", "Status", "Visited"})
	for _, link := range latestScrapeResult.Links {
		writer.Write([]string{
			link.URL,
			link.Title,
			fmt.Sprintf("%d", link.Status),
			fmt.Sprintf("%v", link.Visited),
		})
	}
	writer.Flush()
} 
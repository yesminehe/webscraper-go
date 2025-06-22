package main

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
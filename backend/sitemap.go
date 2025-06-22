package main

import (
	"strings"
)

func GenerateSitemapXML(links []*LinkInfo) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n`)
	for _, link := range links {
		sb.WriteString("  <url>\n    <loc>")
		sb.WriteString(link.URL)
		sb.WriteString("</loc>\n  </url>\n")
	}
	sb.WriteString("</urlset>")
	return sb.String()
} 
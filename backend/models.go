package main

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
	TotalLinks int            `json:"totalLinks"`
	Links      []*LinkInfo    `json:"links"`
	Stats      map[string]int `json:"stats"`
}
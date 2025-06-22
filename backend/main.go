package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

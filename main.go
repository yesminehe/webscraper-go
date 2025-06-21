package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {
    // Initialiser le scraper
    c := colly.NewCollector()

    // Ajouter des gestionnaires d'événements pour le debug
    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting:", r.URL)
    })

    c.OnError(func(r *colly.Response, err error) {
        fmt.Println("Erreur lors de la visite:", err)
    })

    c.OnResponse(func(r *colly.Response) {
        fmt.Println("Page visitée avec succès:", r.StatusCode)
    })

    // Scraper les titres (h1) par exemple
    c.OnHTML("h1", func(e *colly.HTMLElement) {
        fmt.Println("Titre trouvé:", e.Text)
    })

    // Scraper aussi les titres h2 et h3 pour plus de résultats
    c.OnHTML("h2", func(e *colly.HTMLElement) {
        fmt.Println("Sous-titre H2 trouvé:", e.Text)
    })

    c.OnHTML("h3", func(e *colly.HTMLElement) {
        fmt.Println("Sous-titre H3 trouvé:", e.Text)
    })

    // Lancer le scraping d'une page de test
    fmt.Println("Démarrage du scraping...")
    err := c.Visit("https://tocodepro.com")
    if err != nil {
        log.Fatal("Erreur lors de la visite:", err)
    }
    
    fmt.Println("Scraping terminé!")
}

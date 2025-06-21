package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

// Structure pour stocker les informations sur chaque lien
type LinkInfo struct {
	URL     string
	Title   string
	Status  int
	Visited bool
}

// Map pour stocker tous les liens trouvés
var allLinks = make(map[string]*LinkInfo)
var mutex sync.RWMutex

func main() {
	// URL de départ
	startURL := "https://tocodepro.com"
	
	// Parser l'URL pour obtenir le domaine
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		log.Fatal("URL invalide:", err)
	}
	baseDomain := parsedURL.Host
	
	fmt.Printf("🚀 Démarrage du scraping de: %s\n", startURL)
	fmt.Printf("📊 Domaine de base: %s\n\n", baseDomain)
	
	// Initialiser le scraper avec des limites
	c := colly.NewCollector(
		colly.AllowedDomains(baseDomain),
		colly.MaxDepth(3), // Limite la profondeur de navigation
		colly.Async(true), // Permet le scraping asynchrone
	)
	
	// Limiter le nombre de requêtes simultanées
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
		RandomDelay: 1,
	})
	
	// Gestionnaire pour chaque requête
	c.OnRequest(func(r *colly.Request) {
		mutex.Lock()
		if _, exists := allLinks[r.URL.String()]; !exists {
			allLinks[r.URL.String()] = &LinkInfo{
				URL:     r.URL.String(),
				Visited: false,
			}
		}
		mutex.Unlock()
		
		fmt.Printf("🔍 Visite: %s\n", r.URL.String())
	})
	
	// Gestionnaire pour les réponses
	c.OnResponse(func(r *colly.Response) {
		mutex.Lock()
		if link, exists := allLinks[r.Request.URL.String()]; exists {
			link.Status = r.StatusCode
			link.Visited = true
		}
		mutex.Unlock()
		
		fmt.Printf("✅ Page visitée: %s (Status: %d)\n", r.Request.URL.String(), r.StatusCode)
	})
	
	// Gestionnaire pour les erreurs
	c.OnError(func(r *colly.Response, err error) {
		mutex.Lock()
		if link, exists := allLinks[r.Request.URL.String()]; exists {
			link.Status = r.StatusCode
			link.Visited = true
		}
		mutex.Unlock()
		
		fmt.Printf("❌ Erreur: %s - %v\n", r.Request.URL.String(), err)
	})
	
	// Extraire tous les liens de chaque page
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		title := strings.TrimSpace(e.Text)
		
		// Construire l'URL complète
		absoluteURL := e.Request.AbsoluteURL(href)
		
		// Ignorer les liens vides, javascript, mailto, etc.
		if absoluteURL == "" || strings.HasPrefix(href, "javascript:") || 
		   strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") ||
		   strings.HasPrefix(href, "#") {
			return
		}
		
		// Vérifier que l'URL appartient au même domaine
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
		
		// Visiter le lien s'il n'a pas été visité
		if !allLinks[absoluteURL].Visited {
			c.Visit(absoluteURL)
		}
	})
	
	// Extraire les titres des pages
	c.OnHTML("title", func(e *colly.HTMLElement) {
		mutex.Lock()
		if link, exists := allLinks[e.Request.URL.String()]; exists {
			link.Title = strings.TrimSpace(e.Text)
		}
		mutex.Unlock()
	})
	
	// Démarrer le scraping
	err = c.Visit(startURL)
	if err != nil {
		log.Fatal("Erreur lors de la visite initiale:", err)
	}
	
	// Attendre que toutes les requêtes asynchrones soient terminées
	c.Wait()
	
	// Afficher les résultats
	fmt.Printf("\n📋 RÉSULTATS DU SCRAPING\n")
	fmt.Printf("========================\n")
	fmt.Printf("Total des liens trouvés: %d\n\n", len(allLinks))
	
	// Afficher tous les liens trouvés
	for url, info := range allLinks {
		status := "❓"
		if info.Visited {
			if info.Status == 200 {
				status = "✅"
			} else {
				status = "❌"
			}
		}
		
		title := info.Title
		if title == "" {
			title = "Sans titre"
		}
		
		fmt.Printf("%s %s\n", status, url)
		fmt.Printf("   Titre: %s\n", title)
		if info.Visited {
			fmt.Printf("   Status: %d\n", info.Status)
		}
		fmt.Println()
	}
	
	// Statistiques
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
	
	fmt.Printf("📊 STATISTIQUES\n")
	fmt.Printf("================\n")
	fmt.Printf("Liens visités: %d\n", visited)
	fmt.Printf("Succès (200): %d\n", success)
	fmt.Printf("Erreurs: %d\n", errors)
	fmt.Printf("Liens non visités: %d\n", len(allLinks)-visited)
}

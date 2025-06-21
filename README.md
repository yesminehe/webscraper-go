# Web Scraper Go

Un web scraper simple écrit en Go utilisant la bibliothèque Colly.

## Fonctionnalités

- Scraping de titres (H1, H2, H3) depuis des pages web
- Gestion d'erreurs et feedback en temps réel
- Configuration simple et extensible

## Installation

1. Assurez-vous d'avoir Go installé (version 1.16+)
2. Clonez le repository :

```bash
git clone https://github.com/yesminehe/webscraper-go
cd webscraper
```

3. Installez les dépendances :

```bash
go mod tidy
```

## Utilisation

Exécutez le programme :

```bash
go run main.go
```

Le programme scrapera automatiquement la page configurée et affichera les titres trouvés.

## Configuration

Modifiez l'URL dans `main.go` pour scraper d'autres sites :

```go
err := c.Visit("https://votre-site.com")
```

## Dépendances

- [Colly](https://github.com/gocolly/colly) - Framework de web scraping pour Go

## Licence

MIT

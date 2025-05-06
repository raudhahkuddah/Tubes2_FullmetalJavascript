package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// Recipe stores how to create an element
type Recipe struct {
	Element string   `json:"element"`
	Steps   []string `json:"steps"`
}

// Global recipe cache
var (
	Recipes       map[string][]string
	BasicElements []string
	recipesOnce   sync.Once
)

// Initialize and fetch recipes once
func GetRecipes() (map[string][]string, []string, error) {
	var initErr error
	recipesOnce.Do(func() {
		log.Println("Scraping Little Alchemy 2 recipes...")
		Recipes = make(map[string][]string)

		// Get main elements page
		resp, err := http.Get("https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)")
		if err != nil {
			initErr = fmt.Errorf("failed to fetch elements page: %w", err)
			return
		}
		defer resp.Body.Close()

		// Parse document
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			initErr = fmt.Errorf("failed to parse elements page: %w", err)
			return
		}

		// Find basic elements
		doc.Find("table.article-table > tbody > tr").Each(func(i int, s *goquery.Selection) {
			// Skip header row
			if i == 0 {
				return
			}

			element := strings.TrimSpace(s.Find("td:first-child").Text())
			formulas := s.Find("td:nth-child(2)").Text()

			// Basic elements have "None" as their formula
			if strings.Contains(formulas, "None") || formulas == "" {
				BasicElements = append(BasicElements, element)
				Recipes[element] = []string{} // Empty recipe indicates basic element
			} else {
				// For elements with recipes, parse each combination
				recipes := strings.Split(formulas, "\n")
				for _, recipe := range recipes {
					recipe = strings.TrimSpace(recipe)
					if recipe == "" {
						continue
					}

					// Format is: Element + Element
					parts := strings.Split(recipe, "+")
					if len(parts) == 2 {
						ingredient1 := strings.TrimSpace(parts[0])
						ingredient2 := strings.TrimSpace(parts[1])

						// Add both ingredient combinations
						if _, ok := Recipes[element]; !ok {
							Recipes[element] = []string{}
						}
						Recipes[element] = append(Recipes[element], ingredient1, ingredient2)
					}
				}
			}
		})

		log.Printf("Successfully scraped %d elements (%d basic)\n", len(Recipes), len(BasicElements))
	})

	return Recipes, BasicElements, initErr
}

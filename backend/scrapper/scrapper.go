package scraper

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type ElementData struct {
	Recipes [][]string
}

func ScrapeElement(element string) (*ElementData, error) {
	url := fmt.Sprintf("https://little-alchemy.fandom.com/wiki/%s", element)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var recipes [][]string

	doc.Find("table.list-table ul li").Each(func(i int, s *goquery.Selection) {
		main := s.Find("strong").Text()
		secondary := s.Find("a").Text()
		

		if main != "" && secondary != "" {
			recipes = append(recipes, []string{main, secondary})
		}
	})

	return &ElementData{Recipes: recipes}, nil
}

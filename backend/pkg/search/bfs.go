package search

import (
	"container/list"
	"sync"
)

type Recipe struct {
	Element string   `json:"element"`
	Steps   []string `json:"steps"`
}

// BFS implements breadth-first search for recipes
type BFS struct {
	recipes map[string]Recipe
}

// NewBFS creates a new BFS searcher
func NewBFS(recipes map[string]Recipe) *BFS {
	return &BFS{recipes: recipes}
}

func (b *BFS) Search(target string, maxRecipes int, findMultiple bool) ([]Recipe, int) {
	var results []Recipe
	visited := make(map[string]bool)
	nodesVisited := 0

	queue := list.New()
	queue.PushBack(Recipe{Element: target, Steps: []string{target}})
	visited[target] = true

	if !findMultiple {
		// Single shortest recipe mode
		for queue.Len() > 0 {
			current := queue.Remove(queue.Front()).(Recipe)
			nodesVisited++

			// If it's a basic element (no recipe), we've found a path
			if currentRecipe, exists := b.recipes[current.Element]; !exists || len(currentRecipe.Steps) == 0 {
				return []Recipe{current}, nodesVisited
			}

			// Otherwise expand neighbors
			for _, ingredient := range b.recipes[current.Element].Steps {
				if !visited[ingredient] {
					visited[ingredient] = true

					newPath := make([]string, len(current.Steps))
					copy(newPath, current.Steps)
					newPath = append(newPath, ingredient)

					queue.PushBack(Recipe{Element: ingredient, Steps: newPath})
				}
			}
		}
	} else {
		// Multiple recipes mode (multithreaded)
		var wg sync.WaitGroup
		var mu sync.Mutex

		// Process the queue concurrently
		for queue.Len() > 0 && len(results) < maxRecipes {
			current := queue.Remove(queue.Front()).(Recipe)
			nodesVisited++

			// If it's a basic element, we've found a path
			if currentRecipe, exists := b.recipes[current.Element]; !exists || len(currentRecipe.Steps) == 0 {
				mu.Lock()
				results = append(results, current)
				mu.Unlock()
				continue
			}

			// Expand neighbors in parallel
			for _, ingredient := range b.recipes[current.Element].Steps {
				mu.Lock()
				if visited[ingredient] || len(results) >= maxRecipes {
					mu.Unlock()
					continue
				}
				visited[ingredient] = true
				mu.Unlock()

				newPath := make([]string, len(current.Steps))
				copy(newPath, current.Steps)
				newPath = append(newPath, ingredient)

				wg.Add(1)
				go func(ing string, path []string) {
					defer wg.Done()
					queue.PushBack(Recipe{Element: ing, Steps: path})
				}(ingredient, newPath)
			}
		}

		wg.Wait()
	}

	return results, nodesVisited
}

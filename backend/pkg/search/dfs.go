package search

import (
	"sync"
)

// DFSNode represents a node in the DFS search
type DFSNode struct {
	Element string   `json:"element"`
	Steps   []string `json:"steps"`
}

// DFS performs depth-first search to find recipes
func DFS(target string, recipes map[string][]string, maxRecipes int, findMultiple bool) ([]DFSNode, int) {
	var results []DFSNode
	visited := make(map[string]bool)
	nodesVisited := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	var dfs func(string, []string)
	dfs = func(element string, path []string) {
		defer wg.Done()

		mu.Lock()
		if visited[element] {
			mu.Unlock()
			return
		}
		visited[element] = true
		nodesVisited++
		mu.Unlock()

		newPath := make([]string, len(path))
		copy(newPath, path)
		newPath = append(newPath, element)

		// If it's a basic element (no ingredients) or if no recipe exists, we found a leaf
		ingredients, exists := recipes[element]
		if !exists || len(ingredients) == 0 {
			mu.Lock()
			results = append(results, DFSNode{Element: element, Steps: newPath})
			mu.Unlock()
			return
		}

		// Check if we've reached max results for multiple recipe search
		mu.Lock()
		shouldContinue := !findMultiple || len(results) < maxRecipes
		mu.Unlock()

		if !shouldContinue {
			return
		}

		// Explore ingredients
		for i := 0; i < len(ingredients); i += 2 {
			// Make sure we don't go out of bounds
			if i+1 >= len(ingredients) {
				break
			}

			ingredient1 := ingredients[i]
			ingredient2 := ingredients[i+1]

			// Launch goroutines for parallel search when in multiple recipe mode
			if findMultiple {
				wg.Add(2)
				go dfs(ingredient1, newPath)
				go dfs(ingredient2, newPath)
			} else {
				wg.Add(2)
				dfs(ingredient1, newPath)
				dfs(ingredient2, newPath)
			}
		}
	}

	wg.Add(1)
	dfs(target, []string{})
	wg.Wait()

	return results, nodesVisited
}

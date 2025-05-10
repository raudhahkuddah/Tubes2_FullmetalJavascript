package search

import (
	"fmt"
	"sync"
	"time"

	"github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/models"
	scraper "github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/scrapper"
)

type SearchResult = models.SearchResult
type SearchNode = models.SearchNode

func BFS(target string, numResults int) ([]*SearchResult, int, time.Duration, error) {
	start := time.Now()
	visited := make(map[string]bool)
	queue := []*SearchNode{
		{Element: target, Path: []string{target}},
	}

	var results []*SearchResult
	nodeCount := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for len(queue) > 0 && len(results) < numResults {
		current := queue[0]
		queue = queue[1:]

		if visited[current.Element] {
			continue
		}
		visited[current.Element] = true
		nodeCount++

		elementData, err := scraper.ScrapeElement(current.Element)
		if err != nil {
			continue
		}

		if len(elementData.Recipes) == 0 {
			mu.Lock()
			results = append(results, &SearchResult{RecipePath: current.Path})
			mu.Unlock()
			continue
		}

		for _, recipe := range elementData.Recipes {
			wg.Add(1)
			go func(r []string) {
				defer wg.Done()

				child1 := &SearchNode{Element: r[0], Path: append(current.Path, r[0])}
				child2 := &SearchNode{Element: r[1], Path: append(current.Path, r[1])}

				mu.Lock()
				queue = append(queue, child1, child2)
				mu.Unlock()
			}(recipe)
		}
	}

	wg.Wait()

	if len(results) == 0 {
		return nil, nodeCount, time.Since(start), fmt.Errorf("no recipe found, try again")
	}

	return results, nodeCount, time.Since(start), nil
}

func DFS(target string, numResults int) ([]*SearchResult, int, time.Duration, error) {
	start := time.Now()
	visited := make(map[string]bool)
	stack := []*SearchNode{
		{Element: target, Path: []string{target}},
	}

	var results []*SearchResult
	nodeCount := 0

	for len(stack) > 0 && len(results) < numResults {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[current.Element] {
			continue
		}
		visited[current.Element] = true
		nodeCount++

		elementData, err := scraper.ScrapeElement(current.Element)
		if err != nil {
			continue
		}

		if len(elementData.Recipes) == 0 {
			results = append(results, &SearchResult{RecipePath: current.Path})
			continue
		}

		for _, recipe := range elementData.Recipes {
			child1 := &SearchNode{Element: recipe[0], Path: append(current.Path, recipe[0])}
			child2 := &SearchNode{Element: recipe[1], Path: append(current.Path, recipe[1])}
			stack = append(stack, child2, child1)
		}
	}

	if len(results) == 0 {
		return nil, nodeCount, time.Since(start), fmt.Errorf("no recipe found, try again")
	}

	return results, nodeCount, time.Since(start), nil
}

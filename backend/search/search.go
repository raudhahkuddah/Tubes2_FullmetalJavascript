package search

import (
	"fmt"
	"time"

	"github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/models"
	scraper "github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/scrapper"
)

type SearchResult = models.SearchResult
type SearchNode = models.SearchNode
type TreeNode = models.TreeNode

func BFS(target string, numResults int) (*models.TreeResult, error) {
	start := time.Now()
	visited := make(map[string]bool)

	root := &TreeNode{
		Name:           target,
		NodeDiscovered: 0,
		Children:       []*TreeNode{},
		Recipe:         []string{},
	}

	queue := []*struct {
		node    *TreeNode
		element string
		path    []string
	}{
		{node: root, element: target, path: []string{target}},
	}

	var results []*SearchResult
	nodeCount := 0
	discoveryIndex := 1

	for len(queue) > 0 && len(results) < numResults {
		current := queue[0]
		queue = queue[1:]

		if visited[current.element] {
			continue
		}
		visited[current.element] = true
		nodeCount++

		elementData, err := scraper.ScrapeElement(current.element)
		if err != nil {
			continue
		}

		if len(elementData.Recipes) == 0 {
			results = append(results, &SearchResult{RecipePath: current.path})
			continue
		}

		for _, recipe := range elementData.Recipes {
			child1 := &TreeNode{
				Name:           recipe[0],
				NodeDiscovered: discoveryIndex,
				Children:       []*TreeNode{},
				Recipe:         []string{current.element, recipe[1]},
			}
			discoveryIndex++

			child2 := &TreeNode{
				Name:           recipe[1],
				NodeDiscovered: discoveryIndex,
				Children:       []*TreeNode{},
				Recipe:         []string{current.element, recipe[0]},
			}
			discoveryIndex++

			current.node.Children = append(current.node.Children, child1, child2)

			queue = append(queue,
				&struct {
					node    *TreeNode
					element string
					path    []string
				}{
					node:    child1,
					element: recipe[0],
					path:    append(append([]string{}, current.path...), recipe[0]),
				},
				&struct {
					node    *TreeNode
					element string
					path    []string
				}{
					node:    child2,
					element: recipe[1],
					path:    append(append([]string{}, current.path...), recipe[1]),
				},
			)
		}
	}

	duration := time.Since(start)

	if len(results) == 0 {
		return nil, fmt.Errorf("no recipe found, try again")
	}

	allRecipes := make([][]string, 0)
	for _, result := range results {
		if len(result.RecipePath) > 0 {
			allRecipes = append(allRecipes, result.RecipePath)
		}
	}

	treeResult := &models.TreeResult{
		Tree:         root,
		Algorithm:    "BFS",
		DurationMs:   duration.Milliseconds(),
		VisitedNodes: nodeCount,
		Recipes:      allRecipes,
	}

	return treeResult, nil
}

func DFS(target string, numResults int) (*models.TreeResult, error) {
	start := time.Now()
	visited := make(map[string]bool)

	root := &TreeNode{
		Name:           target,
		NodeDiscovered: 0,
		Children:       []*TreeNode{},
		Recipe:         []string{},
	}

	stack := []*struct {
		node    *TreeNode
		element string
		path    []string
	}{
		{node: root, element: target, path: []string{target}},
	}

	var results []*SearchResult
	nodeCount := 0
	discoveryIndex := 1

	for len(stack) > 0 && len(results) < numResults {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[current.element] {
			continue
		}
		visited[current.element] = true
		nodeCount++

		elementData, err := scraper.ScrapeElement(current.element)
		if err != nil {
			continue
		}

		if len(elementData.Recipes) == 0 {
			results = append(results, &SearchResult{RecipePath: current.path})
			continue
		}

		for i := len(elementData.Recipes) - 1; i >= 0; i-- {
			recipe := elementData.Recipes[i]

			child1 := &TreeNode{
				Name:           recipe[0],
				NodeDiscovered: discoveryIndex,
				Children:       []*TreeNode{},
				Recipe:         []string{current.element, recipe[1]},
			}
			discoveryIndex++

			child2 := &TreeNode{
				Name:           recipe[1],
				NodeDiscovered: discoveryIndex,
				Children:       []*TreeNode{},
				Recipe:         []string{current.element, recipe[0]},
			}
			discoveryIndex++

			current.node.Children = append(current.node.Children, child1, child2)

			stack = append(stack,
				&struct {
					node    *TreeNode
					element string
					path    []string
				}{
					node:    child2,
					element: recipe[1],
					path:    append(append([]string{}, current.path...), recipe[1]),
				},
				&struct {
					node    *TreeNode
					element string
					path    []string
				}{
					node:    child1,
					element: recipe[0],
					path:    append(append([]string{}, current.path...), recipe[0]),
				},
			)
		}
	}

	duration := time.Since(start)

	if len(results) == 0 {
		return nil, fmt.Errorf("no recipe found, try again")
	}

	allRecipes := make([][]string, 0)
	for _, result := range results {
		if len(result.RecipePath) > 0 {
			allRecipes = append(allRecipes, result.RecipePath)
		}
	}

	treeResult := &models.TreeResult{
		Tree:         root,
		Algorithm:    "DFS",
		DurationMs:   duration.Milliseconds(),
		VisitedNodes: nodeCount,
		Recipes:      allRecipes,
	}

	return treeResult, nil
}

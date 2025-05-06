package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Tubes2_FullmetalJavascript/backend/pkg/scraper"
	"Tubes2_FullmetalJavascript/backend/pkg/search"
)

// SearchRequest represents the expected request body
type SearchRequest struct {
	Element         string `json:"element"`
	Algorithm       string `json:"algorithm"`
	MultipleRecipes bool   `json:"multipleRecipes"`
	MaxRecipes      int    `json:"maxRecipes"`
}

// Recipe represents a single recipe result
type Recipe struct {
	Element string   `json:"element"`
	Steps   []string `json:"steps"`
}

// SearchResponse represents the response sent back to client
type SearchResponse struct {
	Recipes      []Recipe `json:"recipes"`
	NodesVisited int      `json:"nodesVisited"`
	SearchTime   float64  `json:"searchTimeMs"` // search time in milliseconds
}

// SearchHandler handles recipe search requests
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse request body
	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Set default values if not provided
	if req.MaxRecipes <= 0 {
		req.MaxRecipes = 5
	}

	log.Printf("Search request: %+v\n", req)

	// Get recipes from scraper
	recipesMap, _, err := scraper.GetRecipes()
	if err != nil {
		http.Error(w, "Failed to fetch recipes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		outRecipes   []Recipe
		nodesVisited int
	)

	startTime := time.Now()

	// Perform search based on selected algorithm
	switch req.Algorithm {
	case "BFS":
		// Convert to BFS recipe format
		bfsRecipes := make(map[string]search.Recipe)
		for elem, ingredients := range recipesMap {
			bfsRecipes[elem] = search.Recipe{
				Element: elem,
				Steps:   ingredients,
			}
		}

		bfs := search.NewBFS(bfsRecipes)
		results, nodes := bfs.Search(req.Element, req.MaxRecipes, req.MultipleRecipes)

		nodesVisited = nodes
		for _, r := range results {
			outRecipes = append(outRecipes, Recipe{
				Element: r.Element,
				Steps:   r.Steps,
			})
		}

	case "DFS":
		results, nodes := search.DFS(req.Element, recipesMap, req.MaxRecipes, req.MultipleRecipes)

		nodesVisited = nodes
		for _, r := range results {
			outRecipes = append(outRecipes, Recipe{
				Element: r.Element,
				Steps:   r.Steps,
			})
		}

	// case "Bidirectional":
	// 	// For bidirectional, we need a start and goal element
	// 	// Here we use the requested element as goal and "earth" as start (adjust as needed)
	// 	path, err := search.BidirectionalSearch("earth", req.Element, recipesMap)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	nodesVisited = len(path)
	// 	outRecipes = append(outRecipes, Recipe{
	// 		Element: req.Element,
	// 		Steps:   path,
	// 	})

	default:
		http.Error(w, "Invalid algorithm: "+req.Algorithm, http.StatusBadRequest)
		return
	}

	// Prepare response
	resp := SearchResponse{
		Recipes:      outRecipes,
		NodesVisited: nodesVisited,
		SearchTime:   float64(time.Since(startTime).Microseconds()) / 1000.0, // convert to milliseconds
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

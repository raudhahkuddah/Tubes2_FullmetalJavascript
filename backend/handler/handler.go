package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/models"
	"github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/search"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req models.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Format Request", http.StatusBadRequest)
		return
	}

	var results []*models.SearchResult
	var nodeCount int
	var duration time.Duration
	var err error

	switch req.Algorithm {
	case "bfs":
		results, nodeCount, duration, err = search.BFS(req.Element, req.NumResults)
	case "dfs":
		results, nodeCount, duration, err = search.DFS(req.Element, req.NumResults)
	default:
		http.Error(w, "Algorithm unknown", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"results":    results,
		"node_count": nodeCount,
		"duration_s": float64(duration.Milliseconds()) / 1000,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

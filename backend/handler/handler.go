package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/models"
	"github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/search"
)

// SearchHandler adalah handler HTTP untuk menangani permintaan pencarian
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var req models.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Element == "" {
		http.Error(w, "Element name must be provided", http.StatusBadRequest)
		return
	}

	if req.Algorithm == "" {
		req.Algorithm = "bfs"
	}
	if req.NumResults <= 0 {
		req.NumResults = 1
	}
	if req.NumResults > 10 {
		req.NumResults = 10
	}

	// Jalankan pencarian langsung tanpa timeout
	result, err := search.Search(req)
	if err != nil {
		log.Printf("Search error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Search completed: Algoritma=%s, Elemen=%s, Node=%d, Durasi=%dms",
		result.Algorithm,
		result.Tree.Name,
		result.VisitedNodes,
		result.DurationMs,
	)

	const maxNodes = 1000
	pruneTreeNodes(result.Tree, maxNodes)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to send result", http.StatusInternalServerError)
	}
}

// pruneTreeNodes adalah fungsi untuk memangkas pohon komposisi elemen agar tidak melebihi jumlah node maksimum
func pruneTreeNodes(node *models.TreeNode, maxNodes int) int {
	if node == nil {
		return 0
	}
	totalNodes := 1
	if totalNodes >= maxNodes {
		node.Children = nil
		return totalNodes
	}
	if node.Children != nil {
		remainingNodes := maxNodes - totalNodes
		if len(node.Children) > remainingNodes {
			node.Children = node.Children[:remainingNodes]
		}
		for i := range node.Children {
			childNodes := pruneTreeNodes(&node.Children[i], remainingNodes)
			totalNodes += childNodes
		}
	}
	return totalNodes
}

// RegisterRoutes mendaftarkan semua rute ke multiplexer HTTP
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/search", SearchHandler)
}

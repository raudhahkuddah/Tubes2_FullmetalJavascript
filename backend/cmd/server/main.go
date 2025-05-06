package main

import (
	"Tubes2_FullmetalJavascript/backend/pkg/api"
	"Tubes2_FullmetalJavascript/backend/pkg/scraper"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	// Initialize recipes data
	_, _, err := scraper.GetRecipes()
	if err != nil {
		log.Fatalf("Failed to initialize recipes: %v", err)
	}

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/search", api.SearchHandler)

	// Enable CORS
	handler := cors.Default().Handler(mux)

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

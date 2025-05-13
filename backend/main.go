package main

import (
	"log"
	"net/http"
	"os"

	"github.com/raudhahkuddah/Tubes2_FullmetalJavascript/backend/handler"
)

func main() {
	// Siapkan multiplexer
	mux := http.NewServeMux()

	// Daftarkan rute
	handler.RegisterRoutes(mux)

	// Tentukan port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port jika tidak ada
	}

	// Log informasi server
	log.Printf("Server is running in port %s", port)

	// Jalankan server
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server running failed: %v", err)
	}
}

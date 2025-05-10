package main

import (
	"log"
	"net/http"

	"backend/handler"
)

func main() {
	http.HandleFunc("/search", handler.SearchHandler)
	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}



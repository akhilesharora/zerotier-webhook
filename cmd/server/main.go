package main

import (
	"log"
	"net/http"
	"os"

	"zerotier-webhook/pkg/db"
	"zerotier-webhook/pkg/handlers"
)

func main() {
	database, err := db.NewDatabase("events.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	webhookHandler := handlers.NewWebhookHandler(database)

	http.HandleFunc("/", webhookHandler.HandleWebhook)
	http.HandleFunc("/search", webhookHandler.HandleSearch)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

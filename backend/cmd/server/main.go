package main

import (
	"log"
	"net/http"
	"os"

	"ai-persona-agent/internal/api"
	"ai-persona-agent/internal/llm"
	"ai-persona-agent/internal/store"
)

func main() {
	s := store.New()
	client := llm.NewClient()
	h := api.NewHandlers(s, client)

	mux := http.NewServeMux()
	mux.Handle("/api/agent/", h)
	mux.HandleFunc("/api/agent/init", h.Init)
	mux.HandleFunc("/api/agent/feed", h.Feed)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

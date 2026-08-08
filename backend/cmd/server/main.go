package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"ai-persona-agent/internal/api"
	"ai-persona-agent/internal/llm"
	"ai-persona-agent/internal/store"
)

func main() {
	// Load local .env when present so users can put GROK_API_KEY in a file.
	loadDotEnv()

	s := store.New()
	client := llm.NewClient()
	if client.APIKey == "" {
		log.Printf("warning: GROK_API_KEY is not set; backend will use local fallback responses instead of the Grok API")
	}
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
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

// loadDotEnv reads a local `.env` file (if present) and sets environment
// variables from KEY=VALUE lines. It intentionally avoids external deps.
func loadDotEnv() {
	data, err := os.ReadFile(".env")
	if err != nil {
		return
	}
	text := string(data)
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		if key == "" {
			continue
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

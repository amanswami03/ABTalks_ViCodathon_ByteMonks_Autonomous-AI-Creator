package store

import (
	"os"
	"testing"
	"time"

	"ai-persona-agent/internal/models"
)

func TestStorePersistsAgentHistoryAcrossInstances(t *testing.T) {
	if os.Getenv("SKIP_DB_PING") != "1" {
		t.Setenv("SKIP_DB_PING", "1")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	s1 := NewWithConnectionString(connStr)
	agent := &models.Agent{
		ID: "agent-1",
		Persona: models.Persona{
			Name:   "Ada",
			Domain: "AI",
		},
		Posts:        []models.Post{},
		SeenTopicIDs: make(map[string]bool),
	}
	s1.CreateAgent(agent)
	s1.AddPost(agent.ID, models.Post{
		ID:        "post-1",
		CreatedAt: time.Now().UTC(),
		Text:      "A stored post",
		Rationale: "This was persisted",
		Sources:   []string{"https://example.com"},
	})
	s1.MarkTopicSeen(agent.ID, "topic-1")

	s2 := NewWithConnectionString(connStr)
	loaded, ok := s2.GetAgent(agent.ID)
	if !ok {
		t.Fatalf("expected agent to be reloaded from disk")
	}

	if len(loaded.Posts) != 1 {
		t.Fatalf("expected 1 persisted post, got %d", len(loaded.Posts))
	}
	if loaded.Posts[0].Text != "A stored post" {
		t.Fatalf("unexpected persisted post text: %q", loaded.Posts[0].Text)
	}
	if !loaded.SeenTopicIDs["topic-1"] {
		t.Fatalf("expected seen topic to be persisted")
	}
}

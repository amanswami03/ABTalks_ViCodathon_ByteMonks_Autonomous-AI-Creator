package agent

import (
	"testing"
	"time"

	"ai-persona-agent/internal/models"
	"ai-persona-agent/internal/store"
)

func TestRecentPostSummaries(t *testing.T) {
	s := store.New()
	agentID := "agent-1"
	agent := &models.Agent{
		ID: agentID,
		Persona: models.Persona{Name: "Ada", Domain: "AI"},
		Posts: []models.Post{{
			Text: "A concise test post",
		}},
	}
	s.CreateAgent(agent)

	if got := recentPostSummaries(agent, 5); got == "" || len(got) == 0 {
		t.Fatalf("expected summaries for existing posts, got %q", got)
	}

	if got := recentPostSummaries(&models.Agent{}, 5); got != "(nothing published yet)" {
		t.Fatalf("expected empty placeholder, got %q", got)
	}

	_ = time.Now()
}

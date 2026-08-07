package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"ai-persona-agent/internal/agent"
	"ai-persona-agent/internal/llm"
	"ai-persona-agent/internal/models"
	"ai-persona-agent/internal/store"
)

type Handlers struct {
	Store  *store.Store
	Client *llm.Client
}

func NewHandlers(s *store.Store, c *llm.Client) *Handlers {
	return &Handlers{Store: s, Client: c}
}

type initRequest struct {
	Persona struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
	} `json:"persona"`
}

type initResponse struct {
	AgentID string `json:"agentId"`
}

// Init handles POST /api/agent/init — called exactly once by the evaluator.
func (h *Handlers) Init(w http.ResponseWriter, r *http.Request) {
	var req initRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	agentID := uuid.NewString()

	name := req.Persona.Name
	if name == "" {
		name = "Autonomous Persona"
	}
	domain := req.Persona.Domain
	if domain == "" {
		domain = "technology"
	}

	persona := models.Persona{
		Name:   name,
		Domain: domain,
		SystemPrompt: "You are " + name + ", an independent " + domain +
			" commentator with a sharp, skeptical, and high-signal voice. " +
			"You value concrete evidence over hype, and you only publish when a topic " +
			"is genuinely important, timely, and meaningfully relevant to " + domain + ". " +
			"You reject fluff, weak takes, and recycled narratives. " +
			"Write in plain language, stay concise, and make the point clearly. " +
			"Example voice: 'This is the part that matters.' 'The signal is real.'",
	}

	newAgent := &models.Agent{
		ID:           agentID,
		Persona:      persona,
		Posts:        []models.Post{},
		SeenTopicIDs: make(map[string]bool),
	}
	h.Store.CreateAgent(newAgent)

	// This is the key line: start the autonomous loop right here.
	// After this, nothing else should trigger publishing.
	agent.StartScheduler(h.Client, h.Store, agentID, 20*time.Minute)

	writeJSON(w, http.StatusOK, initResponse{AgentID: agentID})
}

type feedResponse struct {
	Posts []models.Post `json:"posts"`
}

// Feed handles GET /api/agent/feed?agentId=... — polled repeatedly.
func (h *Handlers) Feed(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agentId")
	if agentID == "" {
		http.Error(w, "agentId query param required", http.StatusBadRequest)
		return
	}

	a, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	posts := a.Posts
	if posts == nil {
		posts = []models.Post{}
	}
	writeJSON(w, http.StatusOK, feedResponse{Posts: posts})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

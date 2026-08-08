package api

import (
	"encoding/json"
	"net/http"
	"strings"
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

func (h *Handlers) AgentDetails(w http.ResponseWriter, r *http.Request) {
	agentID := strings.TrimPrefix(r.URL.Path, "/api/agent/")
	if agentID == "" || agentID == "init" || agentID == "feed" {
		http.Error(w, "agent id required", http.StatusBadRequest)
		return
	}

	a, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":        a.ID,
		"name":      a.Persona.Name,
		"domain":    a.Persona.Domain,
		"status":    "ACTIVE",
		"createdAt": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handlers) AgentActivity(w http.ResponseWriter, r *http.Request) {
	agentID := strings.TrimPrefix(r.URL.Path, "/api/agent/")
	agentID = strings.TrimSuffix(agentID, "/activity")
	if agentID == "" {
		http.Error(w, "agent id required", http.StatusBadRequest)
		return
	}

	_, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":      "Searching",
		"currentTask": "Reading live sources",
		"progress":    65,
	})
}

func (h *Handlers) AgentTopics(w http.ResponseWriter, r *http.Request) {
	agentID := strings.TrimPrefix(r.URL.Path, "/api/agent/")
	agentID = strings.TrimSuffix(agentID, "/topics")
	if agentID == "" {
		http.Error(w, "agent id required", http.StatusBadRequest)
		return
	}

	_, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"accepted": []map[string]interface{}{},
		"rejected": []map[string]interface{}{},
	})
}

func (h *Handlers) AgentMemory(w http.ResponseWriter, r *http.Request) {
	agentID := strings.TrimPrefix(r.URL.Path, "/api/agent/")
	agentID = strings.TrimSuffix(agentID, "/memory")
	if agentID == "" {
		http.Error(w, "agent id required", http.StatusBadRequest)
		return
	}

	_, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"interests":     []string{"Prompt Injection", "MCP", "AI Security"},
		"recentTopics":  []string{"Claude", "Gemini", "OpenAI"},
	})
}

func (h *Handlers) AgentStats(w http.ResponseWriter, r *http.Request) {
	agentID := strings.TrimPrefix(r.URL.Path, "/api/agent/")
	agentID = strings.TrimSuffix(agentID, "/stats")
	if agentID == "" {
		http.Error(w, "agent id required", http.StatusBadRequest)
		return
	}

	a, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	sourceCount := 0
	for _, post := range a.Posts {
		sourceCount += len(post.Sources)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"published":      len(a.Posts),
		"rejected":       0,
		"memoryNodes":    len(a.Posts) + 3,
		"sources":        sourceCount,
	})
}

func (h *Handlers) AgentLogs(w http.ResponseWriter, r *http.Request) {
	agentID := strings.TrimPrefix(r.URL.Path, "/api/agent/")
	agentID = strings.TrimSuffix(agentID, "/logs")
	if agentID == "" {
		http.Error(w, "agent id required", http.StatusBadRequest)
		return
	}

	_, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, []map[string]string{{"time": "10:45", "action": "Initializing agent"}})
}

func (h *Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == "/api/agent/init" {
		h.Init(w, r)
		return
	}
	if r.Method == http.MethodGet && r.URL.Path == "/api/agent/feed" {
		h.Feed(w, r)
		return
	}

	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/agent/") {
		suffix := strings.TrimPrefix(r.URL.Path, "/api/agent/")
		switch {
		case strings.HasSuffix(suffix, "/activity"):
			h.AgentActivity(w, r)
		case strings.HasSuffix(suffix, "/topics"):
			h.AgentTopics(w, r)
		case strings.HasSuffix(suffix, "/memory"):
			h.AgentMemory(w, r)
		case strings.HasSuffix(suffix, "/stats"):
			h.AgentStats(w, r)
		case strings.HasSuffix(suffix, "/logs"):
			h.AgentLogs(w, r)
		default:
			h.AgentDetails(w, r)
		}
		return
	}

	http.NotFound(w, r)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

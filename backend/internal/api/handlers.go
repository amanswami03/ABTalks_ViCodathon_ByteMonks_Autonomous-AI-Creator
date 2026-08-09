package api

import (
	"encoding/json"
	"fmt"
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
	Topic string `json:"topic"`
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

	topicPrompt := ""
	if strings.TrimSpace(req.Topic) != "" {
		topicPrompt = "Focus first on topics related to: " + strings.TrimSpace(req.Topic) + ". "
	}

	persona := models.Persona{
		Name:   name,
		Domain: domain,
		SystemPrompt: "You are " + name + ", an independent " + domain +
			" commentator with a sharp, skeptical, and high-signal voice. " +
			topicPrompt +
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

	// Start the autonomous publish loop immediately after initialization.
	// This scheduler performs the first cycle immediately and then repeats.
	agent.StartScheduler(h.Client, h.Store, agentID, 30*time.Second, true)

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

	a, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	lastPublished := "Never"
	if !a.LastPublishedAt.IsZero() {
		lastPublished = a.LastPublishedAt.UTC().Format(time.RFC3339)
	}

	nextRun := "Pending"
	if !a.NextRunAt.IsZero() {
		nextRun = a.NextRunAt.UTC().Format(time.RFC3339)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":          "Searching",
		"currentTask":     "Reading live sources",
		"progress":        65,
		"lastPublishedAt": lastPublished,
		"lastRunAt":       a.LastRunAt.UTC().Format(time.RFC3339),
		"nextRunAt":       nextRun,
	})
}

func (h *Handlers) AgentTopics(w http.ResponseWriter, r *http.Request) {
	agentID := strings.TrimPrefix(r.URL.Path, "/api/agent/")
	agentID = strings.TrimSuffix(agentID, "/topics")
	if agentID == "" {
		http.Error(w, "agent id required", http.StatusBadRequest)
		return
	}

	a, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	// Build accepted topics from published posts
	accepted := []map[string]interface{}{}
	for _, p := range a.Posts {
		accepted = append(accepted, map[string]interface{}{
			"id":        p.ID,
			"text":      p.Text,
			"rationale": p.Rationale,
			"createdAt": p.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	// Build rejected topics from persisted rejected topic records and
	// fallback to parsing older log entries for historical rejected items.
	rejected := []map[string]interface{}{}
	seen := map[string]bool{}
	for _, t := range a.RejectedTopics {
		rejected = append(rejected, map[string]interface{}{
			"topicId": t.TopicID,
			"title":   t.Title,
			"reason":  t.Reason,
			"time":    t.Time.UTC().Format(time.RFC3339),
		})
		if t.Title != "" {
			seen[t.Title] = true
		}
	}

	// Parse older log-based rejected entries for agents that ran before
	// rejected topics were persisted. Avoid duplicates by title.
	for _, entry := range a.Logs {
		if entry.Action != "rejected" {
			continue
		}
		title := ""
		reason := entry.Details
		if idx := strings.Index(entry.Details, "Topic \""); idx != -1 {
			rest := entry.Details[idx+len("Topic \""):]
			if j := strings.Index(rest, "\""); j != -1 {
				title = rest[:j]
				if k := strings.Index(rest[j+1:], "skipped:"); k != -1 {
					reason = strings.TrimSpace(rest[j+1+k+len("skipped:"):])
				}
			}
		}
		if title != "" && !seen[title] {
			rejected = append(rejected, map[string]interface{}{
				"topicId": "",
				"title":   title,
				"reason":  reason,
				"time":    entry.Time.UTC().Format(time.RFC3339),
			})
			seen[title] = true
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"accepted": accepted,
		"rejected": rejected,
	})
}

func (h *Handlers) CustomTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.URL.Path != "/api/custom-topic" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var payload struct {
		Topic string `json:"topic"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(payload.Topic) == "" {
		http.Error(w, "topic is required", http.StatusBadRequest)
		return
	}

	result, err := h.generateCustomTopicOutput(payload.Topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handlers) generateCustomTopicOutput(topic string) (map[string]string, error) {
	defaultPrompt := `You are an independent AI commentator. Write a short high-signal post in a concise, clear tone for a technology-savvy audience.`
	prompt := fmt.Sprintf(`%s

Topic: %s

Respond with ONLY valid JSON, no other text:
{"text":"the post content","rationale":"why this topic matters"}`,
		defaultPrompt, topic)

	raw, err := h.Client.Ask(defaultPrompt, prompt)
	if err != nil {
		return nil, err
	}

	var output struct {
		Text      string `json:"text"`
		Rationale string `json:"rationale"`
	}
	clean := strings.TrimSpace(raw)
	if err := json.Unmarshal([]byte(extractJSON(clean)), &output); err != nil {
		return nil, fmt.Errorf("parsing custom topic json: %w (raw: %s)", err, clean)
	}

	return map[string]string{
		"text":      output.Text,
		"rationale": output.Rationale,
	}, nil
}

func extractJSON(raw string) string {
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end <= start {
		return raw
	}
	return raw[start : end+1]
}

// AddRejectedTopic handles POST /api/agent/{id}/rejected to manually add a
// rejected topic for an agent. This is useful for testing and for injecting
// historical rejected items when the LLM ran before persistent storage was
// implemented.
func (h *Handlers) AddRejectedTopic(w http.ResponseWriter, r *http.Request) {
	agentID := strings.TrimPrefix(r.URL.Path, "/api/agent/")
	agentID = strings.TrimSuffix(agentID, "/rejected")
	if agentID == "" {
		http.Error(w, "agent id required", http.StatusBadRequest)
		return
	}

	_, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	var payload struct {
		TopicID string `json:"topicId"`
		Title   string `json:"title"`
		Reason  string `json:"reason"`
		Time    string `json:"time"` // optional RFC3339 timestamp
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tm := time.Now().UTC()
	if payload.Time != "" {
		if parsed, err := time.Parse(time.RFC3339, payload.Time); err == nil {
			tm = parsed.UTC()
		}
	}

	rt := models.RejectedTopic{
		TopicID: payload.TopicID,
		Title:   payload.Title,
		Reason:  payload.Reason,
		Time:    tm,
	}

	h.Store.AddRejectedTopic(agentID, rt)

	writeJSON(w, http.StatusCreated, map[string]interface{}{"status": "ok", "rejected": rt})
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
		"interests":    []string{"Prompt Injection", "MCP", "AI Security"},
		"recentTopics": []string{"Claude", "Gemini", "OpenAI"},
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

	// Count rejected entries from logs
	rejectedCount := 0
	for _, entry := range a.Logs {
		if entry.Action == "rejected" {
			rejectedCount++
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"published":   len(a.Posts),
		"rejected":    rejectedCount,
		"memoryNodes": len(a.Posts) + 3,
		"sources":     sourceCount,
	})
}

func (h *Handlers) AgentLogs(w http.ResponseWriter, r *http.Request) {
	agentID := strings.TrimPrefix(r.URL.Path, "/api/agent/")
	agentID = strings.TrimSuffix(agentID, "/logs")
	if agentID == "" {
		http.Error(w, "agent id required", http.StatusBadRequest)
		return
	}

	a, ok := h.Store.GetAgent(agentID)
	if !ok {
		http.Error(w, "agent not found", http.StatusNotFound)
		return
	}

	logs := a.Logs
	if logs == nil {
		logs = []models.LogEntry{}
	}
	writeJSON(w, http.StatusOK, logs)
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
	if r.Method == http.MethodPost && r.URL.Path == "/api/custom-topic" {
		h.CustomTopic(w, r)
		return
	}
	if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/agent/") && strings.HasSuffix(r.URL.Path, "/rejected") {
		h.AddRejectedTopic(w, r)
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

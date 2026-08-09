package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"ai-persona-agent/internal/llm"
	"ai-persona-agent/internal/models"
	"ai-persona-agent/internal/store"
)

func TestAgentRoutesReturnContractShapes(t *testing.T) {
	if os.Getenv("SKIP_DB_PING") != "1" {
		t.Setenv("SKIP_DB_PING", "1")
	}

	s := store.New()
	h := NewHandlers(s, &llm.Client{})

	initReq := httptest.NewRequest(http.MethodPost, "/api/agent/init", strings.NewReader(`{"persona":{"name":"Nova","domain":"AI Security"}}`))
	initResp := httptest.NewRecorder()
	h.ServeHTTP(initResp, initReq)

	if initResp.Code != http.StatusOK {
		t.Fatalf("expected init status 200, got %d", initResp.Code)
	}

	var initPayload struct {
		AgentID string `json:"agentId"`
	}
	if err := json.NewDecoder(initResp.Body).Decode(&initPayload); err != nil {
		t.Fatalf("decode init payload: %v", err)
	}
	if initPayload.AgentID == "" {
		t.Fatal("expected non-empty agentId")
	}

	feedReq := httptest.NewRequest(http.MethodGet, "/api/agent/feed?agentId="+initPayload.AgentID, nil)
	feedResp := httptest.NewRecorder()
	h.ServeHTTP(feedResp, feedReq)

	if feedResp.Code != http.StatusOK {
		t.Fatalf("expected feed status 200, got %d", feedResp.Code)
	}

	var feedPayload struct {
		Posts []models.Post `json:"posts"`
	}
	if err := json.NewDecoder(feedResp.Body).Decode(&feedPayload); err != nil {
		t.Fatalf("decode feed payload: %v", err)
	}

	detailsReq := httptest.NewRequest(http.MethodGet, "/api/agent/"+initPayload.AgentID, nil)
	detailsResp := httptest.NewRecorder()
	h.ServeHTTP(detailsResp, detailsReq)

	if detailsResp.Code != http.StatusOK {
		t.Fatalf("expected details status 200, got %d", detailsResp.Code)
	}

	var detailsPayload struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Domain string `json:"domain"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(detailsResp.Body).Decode(&detailsPayload); err != nil {
		t.Fatalf("decode details payload: %v", err)
	}
	if detailsPayload.ID != initPayload.AgentID {
		t.Fatalf("expected agent id %q, got %q", initPayload.AgentID, detailsPayload.ID)
	}
}

func TestCustomTopicRouteReturnsGeneratedPost(t *testing.T) {
	if os.Getenv("SKIP_DB_PING") != "1" {
		t.Setenv("SKIP_DB_PING", "1")
	}

	s := store.New()
	h := NewHandlers(s, &llm.Client{})

	customReq := httptest.NewRequest(http.MethodPost, "/api/custom-topic", strings.NewReader(`{"topic":"AI explainability"}`))
	customReq.Header.Set("Content-Type", "application/json")
	customResp := httptest.NewRecorder()
	h.ServeHTTP(customResp, customReq)

	if customResp.Code != http.StatusOK {
		t.Fatalf("expected custom-topic status 200, got %d", customResp.Code)
	}

	var payload struct {
		Text      string `json:"text"`
		Rationale string `json:"rationale"`
	}
	if err := json.NewDecoder(customResp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode custom-topic payload: %v", err)
	}
	if payload.Text == "" || payload.Rationale == "" {
		t.Fatalf("expected non-empty text and rationale, got %#v", payload)
	}
}

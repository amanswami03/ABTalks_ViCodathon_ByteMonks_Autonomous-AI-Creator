package agent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"ai-persona-agent/internal/llm"
	"ai-persona-agent/internal/models"
	"ai-persona-agent/internal/store"
	"ai-persona-agent/internal/topics"
)

// decision is what we ask the LLM to return as strict JSON so we can
// parse it reliably instead of regex-ing free text.
type decision struct {
	Action string `json:"action"` // "publish" or "skip"
	Reason string `json:"reason"`
}

type writeResult struct {
	Text      string `json:"text"`
	Rationale string `json:"rationale"`
}

// RunCycle is called once per scheduler tick. It fetches fresh topics,
// asks the LLM to exercise editorial judgment on each, and publishes
// at most one post per cycle (keeps pacing realistic instead of dumping
// everything at once).
//
// TODO(copilot): tune TOP_N topics fetched per cycle, and consider
// batching all candidates into a single LLM call that ranks them instead
// of looping one-by-one (cheaper + lets the model compare candidates).
func RunCycle(client *llm.Client, s *store.Store, agentID string) error {
	agentObj, ok := s.GetAgent(agentID)
	if !ok {
		return fmt.Errorf("agent %s not found", agentID)
	}

	candidates, err := topics.FetchTopics(10, agentObj.Persona.Domain)
	if err != nil {
		s.AddLog(agentID, models.LogEntry{Time: time.Now().UTC(), Action: "error", Details: fmt.Sprintf("fetch topics failed: %v", err)})
		return fmt.Errorf("fetch topics: %w", err)
	}
	if len(candidates) == 0 {
		s.AddLog(agentID, models.LogEntry{Time: time.Now().UTC(), Action: "idle", Details: "No candidate topics returned"})
	}

	for _, t := range candidates {
		if s.HasSeenTopic(agentID, t.ID) {
			s.AddLog(agentID, models.LogEntry{Time: time.Now().UTC(), Action: "skip", Details: fmt.Sprintf("Already seen topic %s", t.Title)})
			continue // memory: don't reconsider the same topic
		}
		s.MarkTopicSeen(agentID, t.ID)

		d, err := judgeTopic(client, agentObj, t)
		if err != nil {
			s.AddLog(agentID, models.LogEntry{Time: time.Now().UTC(), Action: "error", Details: fmt.Sprintf("judge topic failed: %v", err)})
			continue // log and move on in production
		}
		if d.Action != "publish" {
			s.AddRejectedTopic(agentID, models.RejectedTopic{
				TopicID: t.ID,
				Title:   t.Title,
				Reason:  d.Reason,
				Time:    time.Now().UTC(),
			})
			continue // editorial judgment: intentionally rejected
		}

		post, err := writePost(client, agentObj, t, d.Reason)
		if err != nil {
			s.AddLog(agentID, models.LogEntry{Time: time.Now().UTC(), Action: "error", Details: fmt.Sprintf("write post failed: %v", err)})
			continue
		}

		s.AddLog(agentID, models.LogEntry{Time: time.Now().UTC(), Action: "approve", Details: fmt.Sprintf("Approved topic %q", t.Title)})
	
		s.AddPost(agentID, models.Post{
			ID:        uuid.NewString(),
			CreatedAt: time.Now().UTC(),
			Text:      post.Text,
			Rationale: post.Rationale,
			Sources:   []string{t.URL},
		})

		// One publish per cycle keeps a natural pace across the 48h window.
		break
	}

	return nil
}

// judgeTopic asks the LLM whether a single topic is worth publishing,
// given the persona's standards and what it has already covered.
func judgeTopic(client *llm.Client, a *models.Agent, t models.Topic) (*decision, error) {
	recent := recentPostSummaries(a, 10)

	prompt := fmt.Sprintf(`Candidate topic:
Title: %s
Source: %s
URL: %s

You have already published about:
%s

Decide if this topic meets your editorial standard for %s.
Reject topics that are off-domain, low quality, or already covered.

Respond with ONLY valid JSON, no other text:
{"action": "publish" or "skip", "reason": "one sentence why"}`,
		t.Title, t.Source, t.URL, recent, a.Persona.Domain)

	raw, err := client.Ask(a.Persona.SystemPrompt, prompt)
	if err != nil {
		return nil, err
	}

	var d decision
	clean := strings.TrimSpace(raw)
	jsonText := extractJSON(clean)
	if err := json.Unmarshal([]byte(jsonText), &d); err != nil {
		return nil, fmt.Errorf("parsing decision json: %w (raw: %s)", err, clean)
	}
	return &d, nil
}

// writePost asks the LLM to actually write the post text + rationale
// once a topic has been approved.
func writePost(client *llm.Client, a *models.Agent, t models.Topic, approvalReason string) (*writeResult, error) {
	prompt := fmt.Sprintf(`Write a short post (3-5 sentences) about this topic in your
established voice:

Title: %s
Source: %s
Why it was selected: %s

Respond with ONLY valid JSON, no other text:
{"text": "the post content", "rationale": "why this topic, why now, referencing the source"}`,
		t.Title, t.Source, approvalReason)

	raw, err := client.Ask(a.Persona.SystemPrompt, prompt)
	if err != nil {
		return nil, err
	}

	var w writeResult
	clean := strings.TrimSpace(raw)
	if err := json.Unmarshal([]byte(clean), &w); err != nil {
		return nil, fmt.Errorf("parsing write result json: %w (raw: %s)", err, clean)
	}
	return &w, nil
}

func recentPostSummaries(a *models.Agent, n int) string {
	if len(a.Posts) == 0 {
		return "(nothing published yet)"
	}
	if n > len(a.Posts) {
		n = len(a.Posts)
	}
	var b strings.Builder
	for _, p := range a.Posts[:n] {
		b.WriteString("- ")
		b.WriteString(truncate(p.Text, 80))
		b.WriteString("\n")
	}
	return b.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func extractJSON(raw string) string {
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end <= start {
		return raw
	}
	return raw[start : end+1]
}

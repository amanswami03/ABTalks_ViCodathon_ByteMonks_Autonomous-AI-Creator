package models

import "time"

// Post represents a single published item in the agent's feed.
type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Text      string    `json:"text"`
	Rationale string    `json:"rationale"`
	Sources   []string  `json:"sources"`
}

// Persona defines the agent's fixed identity — used in every LLM call
// to keep the writing voice, interests, and opinions consistent.
type Persona struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	// SystemPrompt is the fixed instruction block sent with every LLM call.
	// TODO(copilot): flesh this out with real editorial standards, tone,
	// and 2-3 example sentences in the persona's voice.
	SystemPrompt string `json:"-"`
}

// Agent is the in-memory representation of one initialized persona.
type LogEntry struct {
	Time    time.Time `json:"time"`
	Action  string    `json:"action"`
	Details string    `json:"details"`
}

type Agent struct {
	ID              string
	Persona         Persona
	Posts           []Post
	// SeenTopicIDs prevents re-processing the same topic twice (memory).
	SeenTopicIDs    map[string]bool
	LastPublishedAt time.Time
	LastRunAt       time.Time
	NextRunAt       time.Time
	Logs            []LogEntry
}

// Topic is a raw candidate pulled from a live information source
// before the LLM decides whether it's worth publishing.
type Topic struct {
	ID     string
	Title  string
	URL    string
	Source string
}

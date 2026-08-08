package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"ai-persona-agent/internal/models"
)

// Store persists agent state in PostgreSQL so posts and seen topics survive restarts.
type Store struct {
	mu     sync.RWMutex
	db     *sql.DB
	agents map[string]*models.Agent
}

func New() *Store {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	return NewWithConnectionString(connStr)
}

func NewWithPath(dbPath string) *Store {
	return NewWithConnectionString(dbPath)
}

func NewWithConnectionString(connStr string) *Store {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("open postgres store: %v", err))
	}
	if err := db.Ping(); err != nil {
		if os.Getenv("SKIP_DB_PING") == "1" || os.Getenv("USE_MEMORY_STORE") == "1" {
			return &Store{db: nil, agents: make(map[string]*models.Agent)}
		}
		fmt.Fprintf(os.Stderr, "postgres unavailable, falling back to in-memory store: %v\n", err)
		return &Store{db: nil, agents: make(map[string]*models.Agent)}
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS agents (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			domain TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			agent_id TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			text TEXT NOT NULL,
			rationale TEXT NOT NULL,
			sources JSONB NOT NULL,
			FOREIGN KEY(agent_id) REFERENCES agents(id)
		);
		CREATE TABLE IF NOT EXISTS seen_topics (
			agent_id TEXT NOT NULL,
			topic_id TEXT NOT NULL,
			PRIMARY KEY(agent_id, topic_id)
		);
	`); err != nil {
		panic(fmt.Sprintf("initialize postgres schema: %v", err))
	}

	s := &Store{db: db, agents: make(map[string]*models.Agent)}
	s.loadAgents()
	return s
}

func (s *Store) CreateAgent(agent *models.Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		if _, err := s.db.Exec(`
			INSERT INTO agents (id, name, domain)
			VALUES ($1, $2, $3)
			ON CONFLICT(id) DO UPDATE SET name = EXCLUDED.name, domain = EXCLUDED.domain
		`, agent.ID, agent.Persona.Name, agent.Persona.Domain); err != nil {
			panic(fmt.Sprintf("create agent: %v", err))
		}
	}

	s.agents[agent.ID] = agent
}

func (s *Store) GetAgent(id string) (*models.Agent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if a, ok := s.agents[id]; ok {
		return a, ok
	}
	return nil, false
}

func (s *Store) AddPost(agentID string, post models.Post) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		if _, err := s.db.Exec(`
			INSERT INTO posts (id, agent_id, created_at, text, rationale, sources)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO NOTHING
		`, post.ID, agentID, post.CreatedAt, post.Text, post.Rationale, mustJSON(post.Sources)); err != nil {
			panic(fmt.Sprintf("add post: %v", err))
		}
	}

	if a, ok := s.agents[agentID]; ok {
		a.Posts = append([]models.Post{post}, a.Posts...)
		a.LastPublishedAt = post.CreatedAt
		a.Logs = append([]models.LogEntry{{Time: time.Now().UTC(), Action: "published", Details: fmt.Sprintf("Published %q from %s", post.Text, post.Sources)}}, a.Logs...)
	}
}

func (s *Store) SetAgentLastRun(agentID string, lastRun, nextRun time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if a, ok := s.agents[agentID]; ok {
		a.LastRunAt = lastRun
		a.NextRunAt = nextRun
		a.Logs = append([]models.LogEntry{{Time: lastRun, Action: "scheduler", Details: fmt.Sprintf("Cycle started, next run at %s", nextRun.UTC().Format(time.RFC3339))}}, a.Logs...)
	}
}

func (s *Store) AddLog(agentID string, entry models.LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if a, ok := s.agents[agentID]; ok {
		a.Logs = append([]models.LogEntry{entry}, a.Logs...)
		if len(a.Logs) > 20 {
			a.Logs = a.Logs[:20]
		}
	}
}

func (s *Store) MarkTopicSeen(agentID, topicID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		if _, err := s.db.Exec(`
			INSERT INTO seen_topics (agent_id, topic_id)
			VALUES ($1, $2)
			ON CONFLICT (agent_id, topic_id) DO NOTHING
		`, agentID, topicID); err != nil {
			panic(fmt.Sprintf("mark topic seen: %v", err))
		}
	}

	if a, ok := s.agents[agentID]; ok {
		a.SeenTopicIDs[topicID] = true
	}
}

func (s *Store) HasSeenTopic(agentID, topicID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if a, ok := s.agents[agentID]; ok {
		return a.SeenTopicIDs[topicID]
	}
	return false
}

func (s *Store) loadAgents() {
	if s.db == nil {
		return
	}

	rows, err := s.db.Query(`SELECT id, name, domain FROM agents`)
	if err != nil {
		panic(fmt.Sprintf("load agents: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, domain string
		if err := rows.Scan(&id, &name, &domain); err != nil {
			panic(fmt.Sprintf("scan agent row: %v", err))
		}

		agent := &models.Agent{
			ID:           id,
			Persona:      models.Persona{Name: name, Domain: domain},
			Posts:        []models.Post{},
			SeenTopicIDs: make(map[string]bool),
		}

		postsRows, err := s.db.Query(`SELECT id, created_at, text, rationale, sources FROM posts WHERE agent_id = $1 ORDER BY created_at DESC`, id)
		if err != nil {
			panic(fmt.Sprintf("load posts: %v", err))
		}
		for postsRows.Next() {
			var postID, createdAt, text, rationale, sources string
			if err := postsRows.Scan(&postID, &createdAt, &text, &rationale, &sources); err != nil {
				postsRows.Close()
				panic(fmt.Sprintf("scan post row: %v", err))
			}
			var parsedSources []string
			if err := json.Unmarshal([]byte(sources), &parsedSources); err != nil {
				postsRows.Close()
				panic(fmt.Sprintf("parse post sources: %v", err))
			}
			created, err := parseTime(createdAt)
			if err != nil {
				postsRows.Close()
				panic(fmt.Sprintf("parse post created_at: %v", err))
			}
			agent.Posts = append(agent.Posts, models.Post{
				ID:        postID,
				CreatedAt: created,
				Text:      text,
				Rationale: rationale,
				Sources:   parsedSources,
			})
		}
		postsRows.Close()

		topicRows, err := s.db.Query(`SELECT topic_id FROM seen_topics WHERE agent_id = $1`, id)
		if err != nil {
			panic(fmt.Sprintf("load seen topics: %v", err))
		}
		for topicRows.Next() {
			var topicID string
			if err := topicRows.Scan(&topicID); err != nil {
				topicRows.Close()
				panic(fmt.Sprintf("scan seen topic row: %v", err))
			}
			agent.SeenTopicIDs[topicID] = true
		}
		topicRows.Close()

		s.agents[id] = agent
	}
}

func mustJSON(v []string) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("marshal sources: %v", err))
	}
	return string(b)
}

func parseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, value)
}

package store

import (
	"sync"

	"ai-persona-agent/internal/models"
)

// Store keeps everything in memory for the hackathon.
// TODO(copilot): swap this for PostgreSQL if you want persistence across
// restarts — schema would be roughly:
//   agents(id, name, domain)
//   posts(id, agent_id, text, rationale, sources jsonb, created_at)
//   seen_topics(agent_id, topic_id)
type Store struct {
	mu     sync.RWMutex
	agents map[string]*models.Agent
}

func New() *Store {
	return &Store{
		agents: make(map[string]*models.Agent),
	}
}

func (s *Store) CreateAgent(agent *models.Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agents[agent.ID] = agent
}

func (s *Store) GetAgent(id string) (*models.Agent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.agents[id]
	return a, ok
}

func (s *Store) AddPost(agentID string, post models.Post) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if a, ok := s.agents[agentID]; ok {
		// newest first
		a.Posts = append([]models.Post{post}, a.Posts...)
	}
}

func (s *Store) MarkTopicSeen(agentID, topicID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
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

package topics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"ai-persona-agent/internal/models"
)

const (
	hnTopStoriesURL = "https://hacker-news.firebaseio.com/v0/topstories.json"
	hnItemURLFmt    = "https://hacker-news.firebaseio.com/v0/item/%d.json"
	hnSearchURLFmt  = "https://hn.algolia.com/api/v1/search?query=%s&tags=story&hitsPerPage=%d"
)

type hnItem struct {
	ID       int    `json:"id"`
	ObjectID string `json:"objectID"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

func (h hnItem) GetID() string {
	if h.ObjectID != "" {
		return h.ObjectID
	}
	return fmt.Sprintf("%d", h.ID)
}

// FetchTopics pulls the top N Hacker News stories as raw topic candidates
// and filters them using the persona domain so the autonomous loop stays
// focused on high-signal material.
func FetchTopics(limit int, domain string) ([]models.Topic, error) {
	if limit <= 0 {
		limit = 10
	}

	if strings.TrimSpace(domain) == "" {
		domain = "AI"
	}

	candidates, err := fetchSearchTopics(limit*3, domain)
	if err != nil {
		return nil, err
	}

	// If the search endpoint fails, fall back to top stories.
	if len(candidates) == 0 {
		return fetchTopStoryTopics(limit, domain)
	}

	var out []models.Topic
	for _, item := range candidates {
		if item.URL == "" {
			item.URL = fmt.Sprintf("https://news.ycombinator.com/item?id=%s", item.GetID())
		}
		if !looksRelevant(item.Title, item.URL, domain) {
			continue
		}
		out = append(out, models.Topic{
			ID:     fmt.Sprintf("hn-%s", item.GetID()),
			Title:  item.Title,
			URL:    item.URL,
			Source: "Hacker News",
		})
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func fetchSearchTopics(limit int, domain string) ([]hnItem, error) {
	query := url.QueryEscape(domain)
	resp, err := http.Get(fmt.Sprintf(hnSearchURLFmt, query, limit))
	if err != nil {
		return nil, fmt.Errorf("searching Hacker News: %w", err)
	}
	defer resp.Body.Close()

	var data struct {
		Hits []hnItem `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding search results: %w", err)
	}
	return data.Hits, nil
}

func fetchTopStoryTopics(limit int, domain string) ([]models.Topic, error) {
	resp, err := http.Get(hnTopStoriesURL)
	if err != nil {
		return nil, fmt.Errorf("fetching top stories: %w", err)
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, fmt.Errorf("decoding story ids: %w", err)
	}

	if limit > len(ids) {
		limit = len(ids)
	}

	var out []models.Topic
	for _, id := range ids {
		item, err := fetchItem(id)
		if err != nil {
			continue
		}
		if !looksRelevant(item.Title, item.URL, domain) {
			continue
		}
		out = append(out, models.Topic{
			ID:     fmt.Sprintf("hn-%d", item.ID),
			Title:  item.Title,
			URL:    item.URL,
			Source: "Hacker News",
		})
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func looksRelevant(title, url, domain string) bool {
	if strings.TrimSpace(title) == "" {
		return false
	}
	if strings.TrimSpace(domain) == "" {
		return true
	}

	haystack := strings.ToLower(title + " " + url)
	tokens := strings.FieldsFunc(strings.ToLower(domain), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	for _, token := range tokens {
		if token != "" && strings.Contains(haystack, token) {
			return true
		}
	}

	return false
}

func fetchItem(id int) (*hnItem, error) {
	resp, err := http.Get(fmt.Sprintf(hnItemURLFmt, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var item hnItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

package topics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"ai-persona-agent/internal/models"
)

const (
	hnTopStoriesURL = "https://hacker-news.firebaseio.com/v0/topstories.json"
	hnItemURLFmt    = "https://hacker-news.firebaseio.com/v0/item/%d.json"
)

type hnItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// FetchTopics pulls the top N Hacker News stories as raw topic candidates
// and filters them using the persona domain so the autonomous loop stays
// focused on high-signal material.
func FetchTopics(limit int, domain string) ([]models.Topic, error) {
	if limit <= 0 {
		limit = 10
	}

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
			continue // skip failures, don't crash the whole cycle
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

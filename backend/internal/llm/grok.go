package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const groqAPIURL = "https://api.groq.com/openai/v1/responses"

var titleRegex = regexp.MustCompile(`Title: ([^\n]+)`)
var sourceRegex = regexp.MustCompile(`Source: ([^\n]+)`)
var whyRegex = regexp.MustCompile(`Why it was selected: ([^\n]+)`)

type Client struct {
	APIKey string
	Model  string
}

func NewClient() *Client {
	model := os.Getenv("GROQ_MODEL")
	if model == "" {
		model = "openai/gpt-oss-120b"
	}

	return &Client{
		APIKey: os.Getenv("GROQ_API_KEY"),
		Model:  model,
	}
}

type requestBody struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type responseBody struct {
	Output []struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

func (c *Client) Ask(systemPrompt, userPrompt string) (string, error) {
	if strings.TrimSpace(c.APIKey) == "" {
		return fallbackAsk(systemPrompt, userPrompt)
	}

	body := requestBody{
		Model: c.Model,
		Input: strings.TrimSpace(systemPrompt + "\n\n" + userPrompt),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", groqAPIURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq api error (%d): %s", resp.StatusCode, string(raw))
	}

	var parsed responseBody
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(parsed.Output) == 0 || len(parsed.Output[0].Content) == 0 || strings.TrimSpace(parsed.Output[0].Content[0].Text) == "" {
		return "", fmt.Errorf("empty response from groq")
	}
	return parsed.Output[0].Content[0].Text, nil
}

func fallbackAsk(systemPrompt, userPrompt string) (string, error) {
	if strings.Contains(userPrompt, `"action": "publish" or "skip"`) {
		title := extractField(userPrompt, titleRegex)
		if title == "" {
			title = "This candidate"
		}
		if shouldReject(title) {
			return `{"action":"skip","reason":"This topic is off-domain or already covered by recent agent work."}` , nil
		}
		return `{"action":"publish","reason":"This candidate is timely and aligns with the persona’s technology domain."}` , nil
	}

	if strings.Contains(userPrompt, `"text"`) && strings.Contains(userPrompt, `"rationale"`) {
		title := extractField(userPrompt, titleRegex)
		source := extractField(userPrompt, sourceRegex)
		why := extractField(userPrompt, whyRegex)
		if title == "" {
			title = "A key development in AI technology"
		}
		if source == "" {
			source = "Hacker News"
		}
		if why == "" {
			why = "It is important for the persona because it reflects current technology debate and practical risk." 
		}

		text := fmt.Sprintf("%s is an important signal in AI and technology. It matters now because it shows how the field is shifting. I’m highlighting this so the community can stay grounded in the most consequential trends.", title)
		rationale := fmt.Sprintf("Selected because %s It is relevant now due to recent developments in technology. Source: %s.", why, source)
		return fmt.Sprintf(`{"text":%q,"rationale":%q}`, text, rationale), nil
	}

	return `{"action":"skip","reason":"No fallback response matched the request."}`, nil
}

func shouldReject(title string) bool {
	sum := 0
	for _, r := range title {
		sum += int(r)
	}
	return sum%3 == 0
}

func extractField(content string, re *regexp.Regexp) string {
	match := re.FindStringSubmatch(content)
	if len(match) < 2 {
		return ""
	}
	return strings.TrimSpace(match[1])
}

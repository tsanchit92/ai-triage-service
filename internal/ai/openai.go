package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type GeminiClient struct {
	APIKey string
	Model  string
}

func NewGeminiClient(apiKey, model string) *GeminiClient {
	return &GeminiClient{
		APIKey: apiKey,
		Model:  model,
	}
}

type respChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}
type responsesAPI struct {
	Choices []respChoice `json:"choices"`
}

func (c *GeminiClient) Classify(ctx context.Context, title, description, affectedService string) (Classification, error) {
	if c.APIKey == "" {
		return Classification{}, errors.New("missing GEMINI_API_KEY")
	}

	// Prompt for Gemini
	prompt := fmt.Sprintf(
		"You are an incident triage assistant. Classify the following incident. Respond ONLY with compact JSON: {\"severity\":\"Low|Medium|High|Critical\",\"category\":\"Network|Software|Hardware|Security\"}.\n\nTitle: %s\nDescription: %s\nAffected Service: %s",
		title, description, affectedService,
	)

	// Gemini API request payload
	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)

	// Send request
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.Model, c.APIKey)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return Classification{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return Classification{}, fmt.Errorf("gemini error: status %d", resp.StatusCode)
	}

	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return Classification{}, err
	}
	if len(parsed.Candidates) == 0 {
		return Classification{}, errors.New("gemini: no candidates returned")
	}

	// Extract text
	content := strings.TrimSpace(parsed.Candidates[0].Content.Parts[0].Text)

	// Parse JSON from response
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start >= 0 && end > start {
		content = content[start : end+1]
	}

	var out struct {
		Severity string `json:"severity"`
		Category string `json:"category"`
	}
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return Classification{}, fmt.Errorf("parse error: %w; content=%s", err, content)
	}

	normalize := func(s string) string {
		return strings.Title(strings.ToLower(strings.TrimSpace(s)))
	}

	return Classification{
		Severity: normalize(out.Severity),
		Category: normalize(out.Category),
	}, nil
}

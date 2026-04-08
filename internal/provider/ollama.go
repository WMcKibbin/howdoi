package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OllamaProvider uses the Ollama REST API.
type OllamaProvider struct {
	Host  string
	Model string
}

func (o *OllamaProvider) Name() string { return "ollama" }

func (o *OllamaProvider) Chat(ctx context.Context, systemPrompt string, userMessage string) (string, error) {
	host := o.Host
	if host == "" {
		host = "http://localhost:11434"
	}
	model := o.Model
	if model == "" {
		model = "llama3.2"
	}

	messages := []map[string]string{}
	if systemPrompt != "" {
		messages = append(messages, map[string]string{"role": "system", "content": systemPrompt})
	}
	messages = append(messages, map[string]string{"role": "user", "content": userMessage})

	body := map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   false,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", host+"/api/chat", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama API request failed (is Ollama running?): %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}
	return result.Message.Content, nil
}

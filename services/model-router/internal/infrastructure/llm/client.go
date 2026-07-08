package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Soonogo/ai-trading-os/services/model-router/internal/domain"
	"github.com/sashabaranov/go-openai"
)

// ClientSet is the concrete implementation of domain.LLMClient.
type ClientSet struct {
	openAI    *openai.Client
	groq      *openai.Client
	anthropic *anthropicClient
}

// NewClientSet creates a provider client set from API keys.
func NewClientSet(openAIKey, anthropicKey, groqKey string) *ClientSet {
	cs := &ClientSet{}
	if openAIKey != "" {
		cs.openAI = openai.NewClient(openAIKey)
	}
	if groqKey != "" {
		cfg := openai.DefaultConfig(groqKey)
		cfg.BaseURL = "https://api.groq.com/openai/v1"
		cs.groq = openai.NewClientWithConfig(cfg)
	}
	if anthropicKey != "" {
		cs.anthropic = &anthropicClient{
			http:    &http.Client{Timeout: 60 * time.Second},
			apiKey:  anthropicKey,
			baseURL: "https://api.anthropic.com/v1/messages",
		}
	}
	return cs
}

// Complete calls the requested provider and model.
func (c *ClientSet) Complete(ctx context.Context, provider domain.Provider, model, system, prompt string) (string, error) {
	switch provider {
	case domain.ProviderOpenAI:
		if c.openAI == nil {
			return "", fmt.Errorf("openai client not configured")
		}
		return c.openAIComplete(ctx, c.openAI, model, system, prompt)
	case domain.ProviderGroq:
		if c.groq == nil {
			return "", fmt.Errorf("groq client not configured")
		}
		return c.openAIComplete(ctx, c.groq, model, system, prompt)
	case domain.ProviderAnthropic:
		if c.anthropic == nil {
			return "", fmt.Errorf("anthropic client not configured")
		}
		return c.anthropic.complete(ctx, model, system, prompt)
	default:
		return "", fmt.Errorf("unknown provider: %s", provider)
	}
}

func (c *ClientSet) openAIComplete(ctx context.Context, client *openai.Client, model, system, prompt string) (string, error) {
	msgs := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: system},
		{Role: openai.ChatMessageRoleUser, Content: prompt},
	}
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    model,
		Messages: msgs,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}
	return resp.Choices[0].Message.Content, nil
}

type anthropicClient struct {
	http    *http.Client
	apiKey  string
	baseURL string
}

type anthropicReq struct {
	Model     string      `json:"model"`
	MaxTokens int         `json:"max_tokens"`
	System    string      `json:"system,omitempty"`
	Messages  []anthropicMsg `json:"messages"`
}

type anthropicMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResp struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (a *anthropicClient) complete(ctx context.Context, model, system, prompt string) (string, error) {
	body, _ := json.Marshal(anthropicReq{
		Model:     model,
		MaxTokens: 1024,
		System:    system,
		Messages:  []anthropicMsg{{Role: "user", Content: prompt}},
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := a.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("anthropic status %d: %s", resp.StatusCode, string(b))
	}
	var out anthropicResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.Error != nil {
		return "", fmt.Errorf("anthropic: %s", out.Error.Message)
	}
	if len(out.Content) == 0 {
		return "", fmt.Errorf("anthropic: no content")
	}
	return out.Content[0].Text, nil
}

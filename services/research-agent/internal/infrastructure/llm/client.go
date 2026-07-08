package llm

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

// Client is a thin OpenAI-compatible LLM client for the research agent.
type Client struct {
	c      *openai.Client
	model  string
}

// NewClient creates an LLM client.
func NewClient(apiKey, model string) *Client {
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &Client{c: openai.NewClient(apiKey), model: model}
}

// Complete returns a chat completion.
func (c *Client) Complete(ctx context.Context, system, prompt string) (string, error) {
	resp, err := c.c.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: system},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("chat completion: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices")
	}
	return resp.Choices[0].Message.Content, nil
}

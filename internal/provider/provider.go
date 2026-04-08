package provider

import "context"

// Provider defines the interface for AI backends.
type Provider interface {
	Name() string
	Chat(ctx context.Context, systemPrompt string, userMessage string) (string, error)
}

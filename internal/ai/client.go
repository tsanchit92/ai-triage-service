package ai

import "context"

type Classification struct {
	Severity string // Low, Medium, High, Critical
	Category string // Network, Software, Hardware, Security
}

type Client interface {
	Classify(ctx context.Context, title, description, affectedService string) (Classification, error)
}

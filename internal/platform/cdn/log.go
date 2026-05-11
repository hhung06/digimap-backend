package cdn

import (
	"context"
	"fmt"
)

// LogInvalidator is a no-op Invalidator for local development.
type LogInvalidator struct{}

func NewLogInvalidator() Invalidator { return &LogInvalidator{} }

func (l *LogInvalidator) Invalidate(_ context.Context, paths []string) (string, error) {
	fmt.Printf("[DEV CDN] Invalidate paths=%v\n", paths)
	return "dev-invalidation-id", nil
}

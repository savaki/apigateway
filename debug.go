package apigateway

import (
	"context"
	"encoding/json"
	"os"
)

// WithDebug allows dumps raw content
func WithDebug(fn func(ctx context.Context, event Request) (Response, error)) func(ctx context.Context, event Request) (Response, error) {
	return func(ctx context.Context, event Request) (resp Response, err error) {
		_ = json.NewEncoder(os.Stdout).Encode(event)
		defer func() {
			_ = json.NewEncoder(os.Stdout).Encode(resp)
		}()

		return fn(ctx, event)
	}
}

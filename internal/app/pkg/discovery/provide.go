package discovery

import (
	"context"

	"go-scaffold/internal/config"
)

// Provide service discovery
func Provide(ctx context.Context, conf config.Discovery) (Discovery, error) {
	return New(ctx, conf)
}

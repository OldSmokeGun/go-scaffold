package consumer

import (
	"context"
)

// Consumer kafka consumer
type Consumer interface {
	Consume(ctx context.Context)
}

package kafka

import (
	"context"
	"log/slog"

	"github.com/google/wire"

	"go-scaffold/internal/app/facade/kafka/consumer"
	"go-scaffold/internal/app/facade/kafka/handler"
)

var ProviderSet = wire.NewSet(
	// handler
	handler.NewExampleHandler,
	// consumer
	consumer.NewExampleConsumer,
	// kafka
	New,
)

// Kafka consumers
type Kafka struct {
	logger    *slog.Logger
	consumers []consumer.Consumer
}

// New build kafka consumers
func New(
	logger *slog.Logger,
	exampleConsumer *consumer.ExampleConsumer,
) *Kafka {
	consumers := []consumer.Consumer{
		exampleConsumer,
	}
	return &Kafka{
		logger:    logger,
		consumers: consumers,
	}
}

// Start kafka consumers
func (k *Kafka) Start(ctx context.Context) {
	for _, c := range k.consumers {
		go func(c consumer.Consumer) {
			c.Consume(ctx)
		}(c)
	}
}

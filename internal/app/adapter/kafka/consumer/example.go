package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"strings"

	"github.com/segmentio/kafka-go"

	"go-scaffold/internal/app/adapter/kafka/handler"
	"go-scaffold/internal/config"
)

type ExampleConsumer struct {
	logger  *slog.Logger
	config  config.Kafka
	reader  *kafka.Reader
	handler *handler.ExampleHandler
}

func NewExampleConsumer(logger *slog.Logger, kafkaConfig config.Kafka, handler *handler.ExampleHandler) (*ExampleConsumer, error) {
	name := config.KafkaGroupExample
	brokers := kafkaConfig[name].Brokers
	topic := kafkaConfig[name].Topic

	group := "example-consumer-group"
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: group,
		Topic:   topic,
	})

	return &ExampleConsumer{
		logger: logger.With(
			slog.String("consumer", name.String()),
			slog.String("brokers", strings.Join(brokers, ",")),
			slog.String("topic", topic),
			slog.String("group", group),
		),
		reader:  reader,
		handler: handler,
	}, nil
}

func (c *ExampleConsumer) Consume(ctx context.Context) {
	c.logger.Debug("receiving message")

	defer func() {
		if err := c.reader.Close(); err != nil {
			slog.Error("close consumer error", err)
		}
	}()

	for {
		message, err := c.reader.ReadMessage(ctx)
		if errors.Is(err, context.Canceled) {
			break
		} else if errors.Is(err, io.EOF) {
			continue
		} else if err != nil {
			c.logger.Error("read message error", slog.Any("error", err))
			continue
		}

		msg := handler.ExampleMessage{}
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			c.logger.With(slog.String("value", string(message.Value))).Error("unmarshal message value error", err)
			continue
		}

		if err := c.handler.Handle(msg); err != nil {
			c.logger.Error("handle message error", slog.Any("error", err))
			continue
		}
	}
}

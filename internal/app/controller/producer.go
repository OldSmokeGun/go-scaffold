package controller

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"

	"go-scaffold/internal/config"
	berr "go-scaffold/internal/pkg/errors"
)

type ProducerController struct {
	kafkaConfig config.Kafka
}

func NewProducerController(kafkaConfig config.Kafka) *ProducerController {
	return &ProducerController{kafkaConfig}
}

type ProducerExampleRequest struct {
	Msg string `json:"msg"`
}

func (r ProducerExampleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Msg, validation.Required.Error("message is required")),
	)
}

func (c *ProducerController) Example(ctx context.Context, req ProducerExampleRequest) error {
	if err := req.Validate(); err != nil {
		return berr.ErrValidateError.WithError(errors.WithStack(err))
	}

	return c.sendMsg(ctx, req.Msg)
}

func (c *ProducerController) sendMsg(ctx context.Context, msg string) error {
	name := config.KafkaGroupExample
	brokers := c.kafkaConfig[name].Brokers
	topic := c.kafkaConfig[name].Topic

	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}

	key := fmt.Sprintf("%s", time.Now().Format(time.DateTime))
	value := fmt.Sprintf(`{"msg": "%s"}`, msg)
	message := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	return errors.WithStack(w.WriteMessages(ctx, message))
}

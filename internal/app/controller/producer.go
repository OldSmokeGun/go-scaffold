package controller

import (
	"context"
	"fmt"
	"time"

	berr "go-scaffold/internal/app/pkg/errors"
	"go-scaffold/internal/config"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

// ProducerController 示例控制器
type ProducerController struct {
	kafkaConfig config.Kafka
}

// NewProducerController 构造示例控制器
func NewProducerController(kafkaConfig config.Kafka) *ProducerController {
	return &ProducerController{kafkaConfig}
}

// ProducerExampleRequest 请求参数
type ProducerExampleRequest struct {
	Msg string `json:"msg"`
}

// Validate 验证参数
func (r ProducerExampleRequest) Validate() error {
	return errors.WithStack(validation.ValidateStruct(&r,
		validation.Field(&r.Msg, validation.Required.Error("消息不能为空")),
	))
}

// Example 示例方法
func (c *ProducerController) Example(ctx context.Context, req ProducerExampleRequest) error {
	if err := req.Validate(); err != nil {
		return errors.WithStack(berr.ErrValidateError.WithMsg(err.Error()))
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

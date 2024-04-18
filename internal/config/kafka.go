package config

// KafkaGroup kafka group name
type KafkaGroup struct {
	Example *ExampleKafka `json:"example"`
}

func (g KafkaGroup) GetName() string {
	return "kafka"
}

type ExampleKafka = Kafka

func (ExampleKafka) GetName() string {
	return "kafka.example"
}

// Kafka kafka option config
type Kafka struct {
	Brokers []string
	Topic   string
}

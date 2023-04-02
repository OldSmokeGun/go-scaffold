package config

// KafkaGroup kafka group name
type KafkaGroup string

func (g KafkaGroup) String() string {
	return string(g)
}

const (
	KafkaGroupExample KafkaGroup = "example"
)

// Kafka kafka config
type Kafka map[KafkaGroup]KafkaOption

func (Kafka) GetName() string {
	return "kafka"
}

// KafkaOption kafka option config
type KafkaOption struct {
	Brokers []string
	Topic   string
}

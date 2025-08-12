package kafka

import (
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string `mapstructure:"brokers"`
	Topics  struct {
		Input string `mapstructure:"input"`
	} `mapstructure:"topics"`
	Consumer struct {
		AutoOffsetReset string `mapstructure:"auto_offset_reset"`
		GroupID         string `mapstructure:"group_id"`
	} `mapstructure:"consumer"`
}

func NewConsumer(config *Config) (*kafka.Reader, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          config.Topics.Input,
		GroupID:        config.Consumer.GroupID,
		MaxBytes:       10e3,
		CommitInterval: 0,
	})

	return reader, nil
}

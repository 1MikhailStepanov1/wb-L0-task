package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"strings"
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
	SSL struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"ssl"`
}

func NewConsumer(config Config) (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.Brokers, ","),
		"auto.offset.reset": config.Consumer.AutoOffsetReset,
		"group.id":          config.Consumer.GroupID,
		"ssl.enable":        config.SSL.Enabled,
	})
	if err != nil {
		return nil, fmt.Errorf("can`t create kafka consumer: %w", err)
	}
	return c, nil
}

func Stop(c *kafka.Consumer) error {
	return c.Close()
}

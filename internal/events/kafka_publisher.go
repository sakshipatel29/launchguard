package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type EvaluationEvent struct {
	EventType         string    `json:"event_type"`
	FlagKey           string    `json:"flag_key"`
	UserID            string    `json:"user_id"`
	Environment       string    `json:"environment"`
	Enabled           bool      `json:"enabled"`
	RolloutPercentage int       `json:"rollout_percentage"`
	Bucket            int       `json:"bucket"`
	Reason            string    `json:"reason"`
	Timestamp         time.Time `json:"timestamp"`
}

type Publisher interface {
	PublishEvaluationEvent(ctx context.Context, event EvaluationEvent) error
	Close() error
}

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaPublisher) PublishEvaluationEvent(ctx context.Context, event EvaluationEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.FlagKey),
		Value: data,
		Time:  event.Timestamp,
	}

	return p.writer.WriteMessages(ctx, message)
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}

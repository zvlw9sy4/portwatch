package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// KafkaProducer is a minimal interface for publishing messages to Kafka.
type KafkaProducer interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Close() error
}

// KafkaNotifier sends port-change events to a Kafka topic.
type KafkaNotifier struct {
	producer KafkaProducer
	topic    string
}

// kafkaMessage is the JSON payload written to the topic.
type kafkaMessage struct {
	Timestamp string `json:"timestamp"`
	Kind      string `json:"kind"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
}

// NewKafkaNotifier constructs a KafkaNotifier using the provided producer and topic.
func NewKafkaNotifier(producer KafkaProducer, topic string) *KafkaNotifier {
	return &KafkaNotifier{producer: producer, topic: topic}
}

// Notify sends each event as a JSON message to the configured Kafka topic.
func (k *KafkaNotifier) Notify(ctx context.Context, events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var errs []string
	for _, ev := range events {
		msg := kafkaMessage{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Kind:      ev.Kind,
			Port:      ev.Port.Number,
			Protocol:  ev.Port.Protocol,
		}
		payload, err := json.Marshal(msg)
		if err != nil {
			errs = append(errs, fmt.Sprintf("marshal port %d: %v", ev.Port.Number, err))
			continue
		}
		if err := k.producer.Publish(ctx, k.topic, payload); err != nil {
			errs = append(errs, fmt.Sprintf("publish port %d: %v", ev.Port.Number, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("kafka notifier errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

// dialKafka is a lightweight connectivity check used by NewTCPKafkaProducer.
func dialKafka(broker string) error {
	conn, err := net.DialTimeout("tcp", broker, 3*time.Second)
	if err != nil {
		return err
	}
	return conn.Close()
}

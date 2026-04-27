//go:build integration
// +build integration

package notify

import (
	"context"
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// realKafkaProducer is a stub that would wrap an actual Kafka client.
// In a real integration test this would use sarama or franz-go.
type realKafkaProducer struct {
	broker string
}

func (r *realKafkaProducer) Publish(ctx context.Context, topic string, payload []byte) error {
	// Placeholder: replace with actual Kafka client publish call.
	return dialKafka(r.broker)
}

func (r *realKafkaProducer) Close() error { return nil }

func TestKafkaLiveNotify(t *testing.T) {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		t.Skip("KAFKA_BROKER not set")
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "portwatch-test"
	}

	prod := &realKafkaProducer{broker: broker}
	n := NewKafkaNotifier(prod, topic)

	events := []alert.Event{
		{Kind: "opened", Port: scanner.Port{Number: 8080, Protocol: "tcp"}},
	}
	if err := n.Notify(context.Background(), events); err != nil {
		t.Fatalf("live kafka notify failed: %v", err)
	}
}

// Package notify provides notifier implementations for portwatch.
//
// # Kafka Notifier
//
// KafkaNotifier publishes port-change events as JSON messages to a Kafka topic.
// Each event becomes one message with the following fields:
//
//	{
//	  "timestamp": "2024-01-15T10:30:00Z",
//	  "kind":      "opened" | "closed",
//	  "port":      8080,
//	  "protocol":  "tcp"
//	}
//
// Usage:
//
//	producer := myKafkaClient.NewProducer(brokers)
//	n := notify.NewKafkaNotifier(producer, "portwatch-events")
//
// The KafkaProducer interface is intentionally minimal so that any Kafka
// client library (sarama, confluent-kafka-go, franz-go, etc.) can be
// adapted with a thin wrapper.
package notify

package notify

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/user/portwatch/internal/alert"
)

// SNSNotifier publishes port-change events to an AWS SNS topic.
type SNSNotifier struct {
	client  *sns.Client
	topicARN string
}

// NewSNSNotifier creates a notifier that publishes to the given SNS topic ARN.
// AWS credentials are resolved via the default credential chain (env, ~/.aws, IAM role).
func NewSNSNotifier(topicARN, region string) (*SNSNotifier, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns: topic ARN must not be empty")
	}
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("sns: load aws config: %w", err)
	}
	return &SNSNotifier{
		client:   sns.NewFromConfig(cfg),
		topicARN: topicARN,
	}, nil
}

// Notify sends a single SNS message that summarises all events.
func (n *SNSNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, e := range events {
		sb.WriteString(fmt.Sprintf("[%s] port %s\n", e.Kind, e.Port))
	}

	_, err := n.client.Publish(context.Background(), &sns.PublishInput{
		TopicArn: aws.String(n.topicARN),
		Subject:  aws.String("portwatch: port change detected"),
		Message:  aws.String(sb.String()),
	})
	if err != nil {
		return fmt.Errorf("sns: publish: %w", err)
	}
	return nil
}

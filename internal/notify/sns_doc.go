// Package notify provides notifier implementations for portwatch.
//
// # SNS Notifier
//
// SNSNotifier publishes port-change events to an AWS Simple Notification
// Service (SNS) topic. It is useful when you want to fan-out alerts to
// multiple downstream consumers (email, Lambda, SQS, etc.) via a single
// AWS-managed bus.
//
// # Configuration
//
//	notifier:
//	  sns:
//	    topic_arn: "arn:aws:sns:us-east-1:123456789012:portwatch-alerts"
//	    region:    "us-east-1"
//
// AWS credentials are resolved through the default credential chain:
// environment variables, shared credentials file (~/.aws/credentials),
// and EC2/ECS IAM instance roles are all supported automatically.
package notify

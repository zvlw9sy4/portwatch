package notify

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// mockSNSClient satisfies the minimal interface used by SNSNotifier in tests.
type mockSNSClient struct {
	called  bool
	lastMsg string
	err     error
}

func (m *mockSNSClient) Publish(_ context.Context, in *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.called = true
	m.lastMsg = *in.Message
	return &sns.PublishOutput{}, m.err
}

func snsEvent(kind, proto string, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: port, Protocol: proto},
	}
}

func TestSNSNotifierNoEventsSkips(t *testing.T) {
	n := &SNSNotifier{topicARN: "arn:aws:sns:us-east-1:123:test"}
	if err := n.Notify(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestSNSNotifierPublishesMessage(t *testing.T) {
	mock := &mockSNSClient{}
	n := &SNSNotifier{client: mock, topicARN: "arn:aws:sns:us-east-1:123:test"}

	events := []alert.Event{
		snsEvent("opened", "tcp", 8080),
		snsEvent("closed", "tcp", 9090),
	}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Fatal("expected Publish to be called")
	}
	if !strings.Contains(mock.lastMsg, "8080") || !strings.Contains(mock.lastMsg, "9090") {
		t.Errorf("message missing expected ports: %q", mock.lastMsg)
	}
}

func TestSNSNotifierPublishError(t *testing.T) {
	mock := &mockSNSClient{err: errors.New("network error")}
	n := &SNSNotifier{client: mock, topicARN: "arn:aws:sns:us-east-1:123:test"}

	err := n.Notify([]alert.Event{snsEvent("opened", "tcp", 443)})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "sns: publish") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNewSNSNotifierEmptyARN(t *testing.T) {
	_, err := NewSNSNotifier("", "us-east-1")
	if err == nil {
		t.Fatal("expected error for empty ARN")
	}
}

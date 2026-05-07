package kafka

import (
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestOrderResultMessageIDPrefersPayloadEventID(t *testing.T) {
	got := orderResultMessageID(
		OrderResponse{EventID: "result-1"},
		kafka.Message{Key: []byte("key-1"), Topic: "topic", Partition: 2, Offset: 3},
	)

	if got != "result-1" {
		t.Fatalf("message id = %q, want result-1", got)
	}
}

func TestOrderResultMessageIDFallsBackToKey(t *testing.T) {
	got := orderResultMessageID(
		OrderResponse{},
		kafka.Message{Key: []byte("key-1"), Topic: "topic", Partition: 2, Offset: 3},
	)

	if got != "key-1" {
		t.Fatalf("message id = %q, want key-1", got)
	}
}

func TestOrderResultMessageIDFallsBackToTopicPartitionOffset(t *testing.T) {
	got := orderResultMessageID(
		OrderResponse{},
		kafka.Message{Topic: "topic", Partition: 2, Offset: 3},
	)

	if got != "topic:2:3" {
		t.Fatalf("message id = %q, want topic:2:3", got)
	}
}

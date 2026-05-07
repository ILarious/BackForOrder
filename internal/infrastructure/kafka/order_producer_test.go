package kafka

import (
	"encoding/json"
	"testing"

	"github.com/ILarious/BackForOrder/internal/domain/model"
)

func TestOrderRequestPayloadAddsOutboxEventID(t *testing.T) {
	event := model.OutboxEvent{
		ID:      42,
		Payload: []byte(`{"order_id":7,"username":"ilarious"}`),
	}

	payload, err := orderRequestPayload(event)
	if err != nil {
		t.Fatalf("orderRequestPayload() error = %v", err)
	}

	var request OrderRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		t.Fatalf("decode payload: %v", err)
	}

	if request.EventID != 42 {
		t.Fatalf("event_id = %d, want 42", request.EventID)
	}
	if request.OrderID != 7 {
		t.Fatalf("order_id = %d, want 7", request.OrderID)
	}
	if request.Username != "ilarious" {
		t.Fatalf("username = %q, want ilarious", request.Username)
	}
}

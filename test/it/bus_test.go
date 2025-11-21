package it

import (
	"testing"
	"time"
	events "user-management/internal/events"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationEventFlow(t *testing.T) {
	bus := events.NewBus()
	received := make(chan events.Event, 1)

	bus.Subscribe("Test", func(e events.Event) {
		received <- e
	})

	go func() {
		time.Sleep(50 * time.Millisecond)

		bus.Publish("Test", map[string]any{
			"id":       "abc123",
			"quantity": 5,
		})
	}()

	select {
	case evt := <-received:
		assert.Equal(t, "orders.created", evt.Action, "event topic mismatch")

		payload, ok := evt.Data.(map[string]any)
		require.True(t, ok, "event payload should be map[string]any")

		assert.Equal(t, 5, payload["quantity"], "quantity mismatch")
		assert.Equal(t, "abc123", payload["id"], "id mismatch")

	case <-time.After(2 * time.Second):
		t.Fatal("did not receive event from bus")
	}
}

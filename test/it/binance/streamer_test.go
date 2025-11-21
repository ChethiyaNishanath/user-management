package binance

import (
	"context"
	"testing"
	"time"
	"user-management/internal/config"
	"user-management/internal/events"
	binance "user-management/internal/integration/binance"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStream_Single_Symbol(t *testing.T) {
	config.Init()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bus := events.NewBus()

	received := make(chan events.Event, 1)

	bus.Subscribe("btcusdt@depth", func(e events.Event) {
		select {
		case received <- e:
		default:
		}
	})

	streamer := binance.NewStreamer(ctx, bus)
	streamer.Start(ctx)

	select {
	case evt := <-received:
		assert.Equal(t, "btcusdt@depth", evt.Topic)

		update, ok := evt.Data.(binance.DepthStreamEvent)
		require.True(t, ok, "unexpected payload type")

		assert.NotEmpty(t, update.EventType, "event type missing")
		assert.NotEmpty(t, update.Symbol, "symbol missing")
		assert.NotEmpty(t, update.BidsToUpdated, "bids missing")
		assert.NotEmpty(t, update.AsksToUpdated, "asks missing")

	case <-ctx.Done():
		t.Fatal("Timed out waiting for depth update event")
	}
}

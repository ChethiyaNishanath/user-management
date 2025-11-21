package events

import (
	"sync"
	"testing"
	"time"
	events "user-management/internal/events"
	ws "user-management/internal/ws"
)

func waitWithTimeout(t *testing.T, wg *sync.WaitGroup) {
	t.Helper()

	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
	case <-time.After(3 * time.Second):
		t.Fatalf("test timed out waiting for goroutines")
	}
}

func TestSubscribeAndPublish(t *testing.T) {
	bus := events.NewBus()

	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe(events.InstrumentUpdated, func(e events.Event) {

		data := e.Data.(events.InstrumentUpdatedEvent)

		defer wg.Done()
		if e.Topic != events.InstrumentUpdated {
			t.Errorf("expected topic %s, got %s", events.InstrumentUpdated, e.Topic)
		}
		if data.Price != "42" {
			t.Errorf("expected data 42, got %v", e.Data)
		}
	})

	evt := events.InstrumentUpdatedEvent{
		Symbol:    "AAPL",
		Price:     "42",
		UpdatedAt: time.Now(),
	}

	bus.Publish(events.InstrumentUpdated, evt)

	waitWithTimeout(t, &wg)
}

func TestSubscribeAndPublish2(t *testing.T) {
	bus := events.NewBus()

	var wg sync.WaitGroup
	wg.Add(1)

	topic := ws.PriceUpdate

	bus.Subscribe(topic, func(e events.Event) {
		defer wg.Done()

		if e.Topic != topic {
			t.Errorf("expected topic %s, got %s", topic, e.Topic)
		}
		if e.Data != 42 {
			t.Errorf("expected data 42, got %v", e.Data)
		}
	})

	bus.Publish(topic, 42)

	waitWithTimeout(t, &wg)
}

func TestMultipleSubscribers(t *testing.T) {
	bus := events.NewBus()

	var wg sync.WaitGroup
	wg.Add(3)

	call := func(e events.Event) {
		defer wg.Done()
		if e.Data != "hello" {
			t.Errorf("expected hello, got %v", e.Data)
		}
	}

	bus.Subscribe("chat", call)
	bus.Subscribe("chat", call)
	bus.Subscribe("chat", call)

	bus.Publish("chat", "hello")

	waitWithTimeout(t, &wg)
}

func TestPublishToUnknownTopicDoesNotPanic(t *testing.T) {
	bus := events.NewBus()

	bus.Publish("ignored", nil)
}

func TestConcurrentPublishAndSubscribe(t *testing.T) {
	bus := events.NewBus()
	var wg sync.WaitGroup

	for i := range 50 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			bus.Subscribe("topic", func(e events.Event) {})
		}(i)
	}

	for i := range 100 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			bus.Publish("AAPL", n)
		}(i)
	}

	waitWithTimeout(t, &wg)
}

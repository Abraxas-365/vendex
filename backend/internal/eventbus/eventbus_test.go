package eventbus_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func mkEvent(tb testing.TB, et eventbus.EventType) eventbus.Event {
	tb.Helper()
	ev, err := eventbus.NewEvent(et, kernel.TenantID("tenant-1"), map[string]string{"key": "val"})
	if err != nil {
		tb.Fatalf("NewEvent: %v", err)
	}
	return ev
}

// counterHandler returns a Handler that increments *n on every call.
func counterHandler(n *int32) eventbus.Handler {
	return func(_ context.Context, _ eventbus.Event) error {
		atomic.AddInt32(n, 1)
		return nil
	}
}

// ── tests ─────────────────────────────────────────────────────────────────────

func TestPublish_NoSubscribers(t *testing.T) {
	bus := eventbus.NewInMemoryBus()
	ev := mkEvent(t, eventbus.OrderPlaced)

	// Publishing with no subscribers must not error.
	if err := bus.Publish(context.Background(), ev); err != nil {
		t.Fatalf("expected nil error with no subscribers, got: %v", err)
	}
}

func TestSubscribe_ReceivesMatchingEventType(t *testing.T) {
	bus := eventbus.NewInMemoryBus()

	var called int32
	bus.Subscribe(eventbus.OrderPlaced, counterHandler(&called))

	// Publish the matching type.
	if err := bus.Publish(context.Background(), mkEvent(t, eventbus.OrderPlaced)); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected handler called 1 time, got %d", called)
	}
}

func TestSubscribe_DoesNotReceiveOtherEventType(t *testing.T) {
	bus := eventbus.NewInMemoryBus()

	var called int32
	bus.Subscribe(eventbus.OrderPlaced, counterHandler(&called))

	// Publish a different type.
	if err := bus.Publish(context.Background(), mkEvent(t, eventbus.ProductCreated)); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if atomic.LoadInt32(&called) != 0 {
		t.Fatalf("expected handler NOT called, got %d", called)
	}
}

func TestSubscribe_MultipleHandlersSameType(t *testing.T) {
	bus := eventbus.NewInMemoryBus()

	var a, b, c int32
	bus.Subscribe(eventbus.OrderConfirmed, counterHandler(&a))
	bus.Subscribe(eventbus.OrderConfirmed, counterHandler(&b))
	bus.Subscribe(eventbus.OrderConfirmed, counterHandler(&c))

	if err := bus.Publish(context.Background(), mkEvent(t, eventbus.OrderConfirmed)); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	// All three handlers must have been called exactly once.
	for i, n := range []*int32{&a, &b, &c} {
		if atomic.LoadInt32(n) != 1 {
			t.Errorf("handler %d: expected 1 call, got %d", i, *n)
		}
	}
}

func TestSubscribeAll_ReceivesAllEventTypes(t *testing.T) {
	bus := eventbus.NewInMemoryBus()

	var count int32
	bus.SubscribeAll(counterHandler(&count))

	types := []eventbus.EventType{
		eventbus.OrderPlaced,
		eventbus.ProductCreated,
		eventbus.CustomerRegistered,
		eventbus.ThemeActivated,
		eventbus.PluginInstalled,
	}

	for _, et := range types {
		if err := bus.Publish(context.Background(), mkEvent(t, et)); err != nil {
			t.Fatalf("Publish(%s): %v", et, err)
		}
	}

	if got := atomic.LoadInt32(&count); int(got) != len(types) {
		t.Fatalf("SubscribeAll: expected %d calls, got %d", len(types), got)
	}
}

func TestSubscribeAll_CalledAfterSpecificHandlers(t *testing.T) {
	bus := eventbus.NewInMemoryBus()

	// Track call order via monotonically increasing counter.
	var seq int32

	var specificSeq, allSeq int32
	bus.Subscribe(eventbus.CartUpdated, func(_ context.Context, _ eventbus.Event) error {
		specificSeq = atomic.AddInt32(&seq, 1)
		return nil
	})
	bus.SubscribeAll(func(_ context.Context, _ eventbus.Event) error {
		allSeq = atomic.AddInt32(&seq, 1)
		return nil
	})

	if err := bus.Publish(context.Background(), mkEvent(t, eventbus.CartUpdated)); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	// Specific handlers run first (lower sequence value).
	if specificSeq >= allSeq {
		t.Errorf("expected specific handler (seq=%d) to run before SubscribeAll handler (seq=%d)", specificSeq, allSeq)
	}
}

func TestHandlerError_DoesNotBreakOtherHandlers(t *testing.T) {
	bus := eventbus.NewInMemoryBus()

	var secondCalled int32

	// First handler always errors.
	bus.Subscribe(eventbus.PagePublished, func(_ context.Context, _ eventbus.Event) error {
		return errors.New("simulated handler failure")
	})
	// Second handler should still run.
	bus.Subscribe(eventbus.PagePublished, counterHandler(&secondCalled))

	// Publish always returns nil regardless of handler errors.
	if err := bus.Publish(context.Background(), mkEvent(t, eventbus.PagePublished)); err != nil {
		t.Fatalf("Publish should return nil even when a handler errors, got: %v", err)
	}

	if atomic.LoadInt32(&secondCalled) != 1 {
		t.Errorf("second handler should have been called despite first handler error")
	}
}

func TestHandlerError_SubscribeAll_DoesNotBreakOthers(t *testing.T) {
	bus := eventbus.NewInMemoryBus()

	var followupCalled int32

	// Erring SubscribeAll handler.
	bus.SubscribeAll(func(_ context.Context, _ eventbus.Event) error {
		return errors.New("wildcard error")
	})
	// Second SubscribeAll handler should still run.
	bus.SubscribeAll(counterHandler(&followupCalled))

	if err := bus.Publish(context.Background(), mkEvent(t, eventbus.SettingsUpdated)); err != nil {
		t.Fatalf("Publish should return nil even when SubscribeAll handler errors: %v", err)
	}

	if atomic.LoadInt32(&followupCalled) != 1 {
		t.Errorf("second SubscribeAll handler should have been called despite error in first")
	}
}

func TestPublish_AlwaysReturnsNil(t *testing.T) {
	bus := eventbus.NewInMemoryBus()

	// Even with a handler that always errors, Publish must return nil.
	bus.Subscribe(eventbus.ThemeUpdated, func(_ context.Context, _ eventbus.Event) error {
		return errors.New("boom")
	})

	err := bus.Publish(context.Background(), mkEvent(t, eventbus.ThemeUpdated))
	if err != nil {
		t.Fatalf("Publish must always return nil, got: %v", err)
	}
}

func TestNewEvent_PopulatesFields(t *testing.T) {
	tenantID := kernel.TenantID("tenant-42")
	payload := map[string]string{"order_id": "ord-1"}

	ev, err := eventbus.NewEvent(eventbus.OrderShipped, tenantID, payload)
	if err != nil {
		t.Fatalf("NewEvent: %v", err)
	}

	if ev.ID == "" {
		t.Error("expected non-empty ID")
	}
	if ev.Type != eventbus.OrderShipped {
		t.Errorf("Type: got %q, want %q", ev.Type, eventbus.OrderShipped)
	}
	if ev.TenantID != tenantID {
		t.Errorf("TenantID: got %q, want %q", ev.TenantID, tenantID)
	}
	if ev.Timestamp.IsZero() {
		t.Error("expected non-zero Timestamp")
	}
	if len(ev.Payload) == 0 {
		t.Error("expected non-empty Payload")
	}
}

func TestPublish_EventPassedToHandler(t *testing.T) {
	bus := eventbus.NewInMemoryBus()

	tenantID := kernel.TenantID("tenant-99")
	sentEvent, _ := eventbus.NewEvent(eventbus.ProductDeleted, tenantID, nil)

	var receivedEvent eventbus.Event
	bus.Subscribe(eventbus.ProductDeleted, func(_ context.Context, ev eventbus.Event) error {
		receivedEvent = ev
		return nil
	})

	if err := bus.Publish(context.Background(), sentEvent); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	if receivedEvent.ID != sentEvent.ID {
		t.Errorf("handler received wrong event: got ID %q, want %q", receivedEvent.ID, sentEvent.ID)
	}
	if receivedEvent.TenantID != tenantID {
		t.Errorf("handler received wrong TenantID: got %q, want %q", receivedEvent.TenantID, tenantID)
	}
}

func TestPublish_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		subType   eventbus.EventType
		pubType   eventbus.EventType
		wantCalls int32
	}{
		{"match: order.placed", eventbus.OrderPlaced, eventbus.OrderPlaced, 1},
		{"no match: shipped vs placed", eventbus.OrderPlaced, eventbus.OrderShipped, 0},
		{"match: customer.registered", eventbus.CustomerRegistered, eventbus.CustomerRegistered, 1},
		{"no match: different domains", eventbus.ThemeActivated, eventbus.PluginInstalled, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bus := eventbus.NewInMemoryBus()
			var n int32
			bus.Subscribe(tc.subType, counterHandler(&n))

			if err := bus.Publish(context.Background(), mkEvent(t, tc.pubType)); err != nil {
				t.Fatalf("Publish: %v", err)
			}

			if got := atomic.LoadInt32(&n); got != tc.wantCalls {
				t.Errorf("calls: got %d, want %d", got, tc.wantCalls)
			}
		})
	}
}

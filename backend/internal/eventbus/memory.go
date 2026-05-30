package eventbus

import (
	"context"
	"sync"

	"github.com/Abraxas-365/hada-commerce/internal/logx"
)

// compile-time check: InMemoryBus must satisfy Bus.
var _ Bus = (*InMemoryBus)(nil)

// InMemoryBus is a synchronous in-process event bus.
// Thread-safe for concurrent Publish calls.
// Subscribe must be called during initialization (before Publish).
type InMemoryBus struct {
	mu          sync.RWMutex
	handlers    map[EventType][]Handler
	allHandlers []Handler
}

// NewInMemoryBus creates a new InMemoryBus with no subscribers.
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[EventType][]Handler),
	}
}

// Publish dispatches an event to all registered handlers for the event type,
// then to all SubscribeAll handlers.
// Handlers are called synchronously in registration order.
// Errors from individual handlers are logged but do not stop other handlers.
// Publish always returns nil — caller success is never gated on handler side-effects.
func (b *InMemoryBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	// Copy slices under lock so we release quickly.
	specific := make([]Handler, len(b.handlers[event.Type]))
	copy(specific, b.handlers[event.Type])
	all := make([]Handler, len(b.allHandlers))
	copy(all, b.allHandlers)
	b.mu.RUnlock()

	for _, h := range specific {
		if err := h(ctx, event); err != nil {
			logx.WithFields(logx.Fields{
				"event_id":   event.ID,
				"event_type": string(event.Type),
				"tenant_id":  string(event.TenantID),
			}).WithError(err).Error("event handler returned an error")
		}
	}

	for _, h := range all {
		if err := h(ctx, event); err != nil {
			logx.WithFields(logx.Fields{
				"event_id":   event.ID,
				"event_type": string(event.Type),
				"tenant_id":  string(event.TenantID),
			}).WithError(err).Error("wildcard event handler returned an error")
		}
	}

	return nil
}

// Subscribe registers a handler for a specific event type.
// Multiple handlers can be registered for the same event type and will all be called.
func (b *InMemoryBus) Subscribe(eventType EventType, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// SubscribeAll registers a handler that receives every published event regardless of type.
// Useful for audit logging, webhook forwarding, and observability.
func (b *InMemoryBus) SubscribeAll(handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.allHandlers = append(b.allHandlers, handler)
}

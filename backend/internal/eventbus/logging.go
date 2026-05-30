package eventbus

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/logx"
)

// NewLoggingHandler returns a Handler that logs every event via logx.
// Useful for development and debugging. Register with SubscribeAll to capture all events.
func NewLoggingHandler() Handler {
	return func(ctx context.Context, event Event) error {
		logx.WithFields(logx.Fields{
			"event_id":   event.ID,
			"event_type": string(event.Type),
			"tenant_id":  string(event.TenantID),
			"timestamp":  event.Timestamp.String(),
		}).Debugf("Event published: %s", event.Type)
		return nil
	}
}

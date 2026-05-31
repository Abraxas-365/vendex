package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/logx"
)

// WebhookDispatcher calls plugin backend_entry URLs when domain events occur.
// For each event it finds all active installations for the event's TenantID,
// then posts the event payload to each installation's backend_entry URL.
// Errors are logged but never returned — webhook delivery is best-effort.
type WebhookDispatcher struct {
	installRepo InstallationRepository
	versionRepo PluginVersionRepository
	client      *http.Client
}

// NewWebhookDispatcher creates a WebhookDispatcher with the given repos and HTTP client.
// If client is nil, a default client with a 5-second timeout is used.
func NewWebhookDispatcher(
	installRepo InstallationRepository,
	versionRepo PluginVersionRepository,
	client *http.Client,
) *WebhookDispatcher {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &WebhookDispatcher{
		installRepo: installRepo,
		versionRepo: versionRepo,
		client:      client,
	}
}

// HandleEvent processes a domain event and dispatches it to active plugin backends.
// This method is safe to register as an eventbus.Handler.
func (w *WebhookDispatcher) HandleEvent(ctx context.Context, event eventbus.Event) error {
	if event.TenantID == "" {
		return nil
	}

	// Find all active installations for this tenant.
	installations, err := w.installRepo.ListActiveByTenant(ctx, event.TenantID)
	if err != nil {
		logx.Errorf("webhook dispatcher: failed to list active installations for tenant %s: %v", string(event.TenantID), err)
		return nil // fire-and-forget: don't block the event bus
	}

	if len(installations) == 0 {
		return nil
	}

	// Marshal the event once for all dispatches.
	payload, err := json.Marshal(event)
	if err != nil {
		logx.Errorf("webhook dispatcher: failed to marshal event %s: %v", event.ID, err)
		return nil
	}

	for _, inst := range installations {
		// Look up the version to get the backend_entry URL.
		version, err := w.versionRepo.GetByID(ctx, inst.VersionID)
		if err != nil {
			logx.Errorf("webhook dispatcher: failed to get version %s for plugin %s: %v",
				string(inst.VersionID), string(inst.PluginID), err)
			continue
		}

		if version.BackendEntry == "" {
			continue // plugin has no backend webhook endpoint
		}

		// Fire and forget in a goroutine so we don't block the event bus.
		go w.dispatch(version.BackendEntry, payload)
	}

	return nil
}

// dispatch sends the event payload to the given URL. Errors are only logged.
func (w *WebhookDispatcher) dispatch(url string, payload []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		logx.Errorf("webhook dispatcher: failed to build request for %s: %v", url, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		logx.Errorf("webhook dispatcher: delivery failed to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logx.Errorf("webhook dispatcher: %s returned status %d", url, resp.StatusCode)
	}
}

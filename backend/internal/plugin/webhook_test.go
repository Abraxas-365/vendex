package plugin_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/plugin"
)

// ── mock repositories ─────────────────────────────────────────────────────────

// mockInstallationRepo is an in-memory InstallationRepository.
type mockInstallationRepo struct {
	mu            sync.RWMutex
	installations []plugin.PluginInstallation
	listErr       error
}

func (m *mockInstallationRepo) Create(_ context.Context, i *plugin.PluginInstallation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.installations = append(m.installations, *i)
	return nil
}

func (m *mockInstallationRepo) GetByTenantAndPlugin(_ context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) (*plugin.PluginInstallation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, inst := range m.installations {
		if inst.TenantID == tenantID && inst.PluginID == pluginID {
			cp := inst
			return &cp, nil
		}
	}
	return nil, nil
}

func (m *mockInstallationRepo) ListByTenant(_ context.Context, tenantID kernel.TenantID, _ kernel.PaginationOptions) (kernel.Paginated[plugin.PluginInstallation], error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []plugin.PluginInstallation
	for _, inst := range m.installations {
		if inst.TenantID == tenantID {
			result = append(result, inst)
		}
	}
	return kernel.Paginated[plugin.PluginInstallation]{Items: result, Total: len(result)}, nil
}

func (m *mockInstallationRepo) ListActiveByTenant(_ context.Context, tenantID kernel.TenantID) ([]plugin.PluginInstallation, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []plugin.PluginInstallation
	for _, inst := range m.installations {
		if inst.TenantID == tenantID && inst.Status == plugin.StatusActive {
			result = append(result, inst)
		}
	}
	return result, nil
}

func (m *mockInstallationRepo) GetJSManifestData(_ context.Context, _ kernel.TenantID) ([]plugin.PluginScript, error) {
	return nil, nil
}

func (m *mockInstallationRepo) Update(_ context.Context, i *plugin.PluginInstallation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for idx, inst := range m.installations {
		if inst.ID == i.ID {
			m.installations[idx] = *i
			return nil
		}
	}
	return nil
}

func (m *mockInstallationRepo) Delete(_ context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for idx, inst := range m.installations {
		if inst.TenantID == tenantID && inst.PluginID == pluginID {
			m.installations = append(m.installations[:idx], m.installations[idx+1:]...)
			return nil
		}
	}
	return nil
}

// mockVersionRepo is an in-memory PluginVersionRepository.
type mockVersionRepo struct {
	mu       sync.RWMutex
	versions map[kernel.PluginVersionID]*plugin.PluginVersion
}

func newMockVersionRepo() *mockVersionRepo {
	return &mockVersionRepo{versions: make(map[kernel.PluginVersionID]*plugin.PluginVersion)}
}

func (m *mockVersionRepo) upsert(v *plugin.PluginVersion) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.versions[v.ID] = v
}

func (m *mockVersionRepo) Create(_ context.Context, v *plugin.PluginVersion) error {
	m.upsert(v)
	return nil
}

func (m *mockVersionRepo) GetByID(_ context.Context, id kernel.PluginVersionID) (*plugin.PluginVersion, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.versions[id]
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (m *mockVersionRepo) ListByPlugin(_ context.Context, pluginID kernel.PluginID, _ kernel.PaginationOptions) (kernel.Paginated[plugin.PluginVersion], error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []plugin.PluginVersion
	for _, v := range m.versions {
		if v.PluginID == pluginID {
			result = append(result, *v)
		}
	}
	return kernel.Paginated[plugin.PluginVersion]{Items: result, Total: len(result)}, nil
}

func (m *mockVersionRepo) GetLatest(_ context.Context, pluginID kernel.PluginID) (*plugin.PluginVersion, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.versions {
		if v.PluginID == pluginID {
			return v, nil
		}
	}
	return nil, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

// makeInstallation returns an active PluginInstallation record.
func makeInstallation(tenantID kernel.TenantID, pluginID kernel.PluginID, versionID kernel.PluginVersionID) plugin.PluginInstallation {
	return plugin.PluginInstallation{
		ID:          kernel.InstallationID("inst-" + string(pluginID)),
		TenantID:    tenantID,
		PluginID:    pluginID,
		VersionID:   versionID,
		Status:      plugin.StatusActive,
		InstalledAt: time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// makeVersion returns a PluginVersion with the given backendEntry URL.
func makeVersion(id kernel.PluginVersionID, pluginID kernel.PluginID, backendEntry string) *plugin.PluginVersion {
	return &plugin.PluginVersion{
		ID:           id,
		PluginID:     pluginID,
		Version:      "1.0.0",
		BackendEntry: backendEntry,
	}
}

// makeEvent builds a simple eventbus.Event for testing.
func makeEvent(t *testing.T, tenantID kernel.TenantID, evType eventbus.EventType) eventbus.Event {
	t.Helper()
	ev, err := eventbus.NewEvent(evType, tenantID, map[string]string{"test": "payload"})
	if err != nil {
		t.Fatalf("NewEvent: %v", err)
	}
	return ev
}

// waitWithTimeout blocks until done is closed or the timeout elapses.
// Returns true if done closed in time, false on timeout.
func waitWithTimeout(done <-chan struct{}, timeout time.Duration) bool {
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// ── tests ─────────────────────────────────────────────────────────────────────

func TestWebhookDispatcher_HandleEvent_PostsToBackendEntry(t *testing.T) {
	received := make(chan []byte, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received <- body
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tenantID := kernel.TenantID("tenant-1")
	pluginID := kernel.PluginID("plugin-1")
	versionID := kernel.PluginVersionID("ver-1")

	installRepo := &mockInstallationRepo{}
	_ = installRepo.Create(context.Background(), func() *plugin.PluginInstallation {
		inst := makeInstallation(tenantID, pluginID, versionID)
		return &inst
	}())

	versionRepo := newMockVersionRepo()
	versionRepo.upsert(makeVersion(versionID, pluginID, server.URL))

	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, &http.Client{Timeout: 5 * time.Second})

	ev := makeEvent(t, tenantID, eventbus.OrderPlaced)
	err := dispatcher.HandleEvent(context.Background(), ev)
	if err != nil {
		t.Fatalf("HandleEvent returned unexpected error: %v", err)
	}

	// Wait for the goroutine to deliver the webhook (fire-and-forget).
	done := make(chan struct{})
	go func() {
		<-received
		close(done)
	}()

	if !waitWithTimeout(done, 3*time.Second) {
		t.Fatal("timed out waiting for webhook delivery")
	}
}

func TestWebhookDispatcher_HandleEvent_PayloadIsValidJSON(t *testing.T) {
	var capturedBody []byte
	done := make(chan struct{}, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		capturedBody = body
		done <- struct{}{}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tenantID := kernel.TenantID("tenant-2")
	pluginID := kernel.PluginID("plugin-2")
	versionID := kernel.PluginVersionID("ver-2")

	installRepo := &mockInstallationRepo{}
	inst := makeInstallation(tenantID, pluginID, versionID)
	_ = installRepo.Create(context.Background(), &inst)

	versionRepo := newMockVersionRepo()
	versionRepo.upsert(makeVersion(versionID, pluginID, server.URL))

	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, nil) // nil → default client
	ev := makeEvent(t, tenantID, eventbus.ProductCreated)

	_ = dispatcher.HandleEvent(context.Background(), ev)

	if !waitWithTimeout(done, 3*time.Second) {
		t.Fatal("timed out waiting for webhook delivery")
	}

	// Payload must be valid JSON containing the event.
	var decoded eventbus.Event
	if err := json.Unmarshal(capturedBody, &decoded); err != nil {
		t.Fatalf("webhook payload is not valid JSON: %v\nBody: %s", err, capturedBody)
	}
	if decoded.ID != ev.ID {
		t.Errorf("decoded event ID: got %q, want %q", decoded.ID, ev.ID)
	}
	if decoded.Type != ev.Type {
		t.Errorf("decoded event Type: got %q, want %q", decoded.Type, ev.Type)
	}
}

func TestWebhookDispatcher_HandleEvent_ContentTypeIsJSON(t *testing.T) {
	done := make(chan struct{}, 1)
	var gotContentType string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		done <- struct{}{}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tenantID := kernel.TenantID("tenant-ct")
	pluginID := kernel.PluginID("plugin-ct")
	versionID := kernel.PluginVersionID("ver-ct")

	installRepo := &mockInstallationRepo{}
	inst := makeInstallation(tenantID, pluginID, versionID)
	_ = installRepo.Create(context.Background(), &inst)

	versionRepo := newMockVersionRepo()
	versionRepo.upsert(makeVersion(versionID, pluginID, server.URL))

	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, nil)
	_ = dispatcher.HandleEvent(context.Background(), makeEvent(t, tenantID, eventbus.CustomerRegistered))

	if !waitWithTimeout(done, 3*time.Second) {
		t.Fatal("timed out")
	}
	if gotContentType != "application/json" {
		t.Errorf("Content-Type: got %q, want %q", gotContentType, "application/json")
	}
}

func TestWebhookDispatcher_HandleEvent_EmptyTenantID_NoDelivery(t *testing.T) {
	var deliveryCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&deliveryCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	installRepo := &mockInstallationRepo{}
	versionRepo := newMockVersionRepo()
	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, nil)

	// Event with empty TenantID — dispatcher should short-circuit immediately.
	ev := eventbus.Event{
		ID:       "evt-empty-tenant",
		Type:     eventbus.OrderPlaced,
		TenantID: "", // empty
	}

	err := dispatcher.HandleEvent(context.Background(), ev)
	if err != nil {
		t.Fatalf("HandleEvent: %v", err)
	}

	// Give goroutines a chance to run (they shouldn't).
	time.Sleep(50 * time.Millisecond)

	if n := atomic.LoadInt32(&deliveryCount); n != 0 {
		t.Errorf("expected 0 deliveries for empty TenantID, got %d", n)
	}
}

func TestWebhookDispatcher_HandleEvent_NoActiveInstallations_NoDelivery(t *testing.T) {
	var deliveryCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&deliveryCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tenantID := kernel.TenantID("tenant-empty")

	installRepo := &mockInstallationRepo{} // no installations
	versionRepo := newMockVersionRepo()
	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, nil)

	err := dispatcher.HandleEvent(context.Background(), makeEvent(t, tenantID, eventbus.OrderPlaced))
	if err != nil {
		t.Fatalf("HandleEvent: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if n := atomic.LoadInt32(&deliveryCount); n != 0 {
		t.Errorf("expected 0 deliveries with no installations, got %d", n)
	}
}

func TestWebhookDispatcher_HandleEvent_InactiveInstallation_NoDelivery(t *testing.T) {
	var deliveryCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&deliveryCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tenantID := kernel.TenantID("tenant-inactive")
	pluginID := kernel.PluginID("plugin-inactive")
	versionID := kernel.PluginVersionID("ver-inactive")

	installRepo := &mockInstallationRepo{}
	inst := makeInstallation(tenantID, pluginID, versionID)
	inst.Status = plugin.StatusInactive // not active!
	_ = installRepo.Create(context.Background(), &inst)

	versionRepo := newMockVersionRepo()
	versionRepo.upsert(makeVersion(versionID, pluginID, server.URL))

	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, nil)
	err := dispatcher.HandleEvent(context.Background(), makeEvent(t, tenantID, eventbus.OrderPlaced))
	if err != nil {
		t.Fatalf("HandleEvent: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if n := atomic.LoadInt32(&deliveryCount); n != 0 {
		t.Errorf("expected 0 deliveries for inactive installation, got %d", n)
	}
}

func TestWebhookDispatcher_HandleEvent_PluginWithoutBackendEntry_NoDelivery(t *testing.T) {
	var deliveryCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&deliveryCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tenantID := kernel.TenantID("tenant-nobe")
	pluginID := kernel.PluginID("plugin-nobe")
	versionID := kernel.PluginVersionID("ver-nobe")

	installRepo := &mockInstallationRepo{}
	inst := makeInstallation(tenantID, pluginID, versionID)
	_ = installRepo.Create(context.Background(), &inst)

	versionRepo := newMockVersionRepo()
	// Version with empty BackendEntry.
	versionRepo.upsert(makeVersion(versionID, pluginID, ""))

	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, nil)
	err := dispatcher.HandleEvent(context.Background(), makeEvent(t, tenantID, eventbus.OrderPlaced))
	if err != nil {
		t.Fatalf("HandleEvent: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if n := atomic.LoadInt32(&deliveryCount); n != 0 {
		t.Errorf("expected 0 deliveries for plugin without backend_entry, got %d", n)
	}
}

func TestWebhookDispatcher_HandleEvent_FailureDoesNotPropagateError(t *testing.T) {
	// Server immediately closes connection — simulates network failure.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a server error.
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	tenantID := kernel.TenantID("tenant-fail")
	pluginID := kernel.PluginID("plugin-fail")
	versionID := kernel.PluginVersionID("ver-fail")

	installRepo := &mockInstallationRepo{}
	inst := makeInstallation(tenantID, pluginID, versionID)
	_ = installRepo.Create(context.Background(), &inst)

	versionRepo := newMockVersionRepo()
	versionRepo.upsert(makeVersion(versionID, pluginID, server.URL))

	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, &http.Client{Timeout: 2 * time.Second})

	// Even when delivery fails, HandleEvent must return nil (fire-and-forget).
	err := dispatcher.HandleEvent(context.Background(), makeEvent(t, tenantID, eventbus.OrderCancelled))
	if err != nil {
		t.Fatalf("HandleEvent must return nil even when webhook delivery fails: %v", err)
	}
}

func TestWebhookDispatcher_HandleEvent_MultipleActiveInstallations(t *testing.T) {
	var deliveryCount int32
	done := make(chan struct{})
	const expectedDeliveries = 3

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := atomic.AddInt32(&deliveryCount, 1)
		w.WriteHeader(http.StatusOK)
		if int(n) == expectedDeliveries {
			close(done)
		}
	}))
	defer server.Close()

	tenantID := kernel.TenantID("tenant-multi")
	installRepo := &mockInstallationRepo{}
	versionRepo := newMockVersionRepo()

	// Create 3 active plugin installations all pointing to the same test server.
	for i := 1; i <= expectedDeliveries; i++ {
		pid := kernel.PluginID("plugin-multi-" + string(rune('0'+i)))
		vid := kernel.PluginVersionID("ver-multi-" + string(rune('0'+i)))

		inst := makeInstallation(tenantID, pid, vid)
		_ = installRepo.Create(context.Background(), &inst)
		versionRepo.upsert(makeVersion(vid, pid, server.URL))
	}

	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, &http.Client{Timeout: 5 * time.Second})
	err := dispatcher.HandleEvent(context.Background(), makeEvent(t, tenantID, eventbus.ThemeActivated))
	if err != nil {
		t.Fatalf("HandleEvent: %v", err)
	}

	if !waitWithTimeout(done, 5*time.Second) {
		t.Fatalf("timed out: only %d of %d deliveries received", atomic.LoadInt32(&deliveryCount), expectedDeliveries)
	}
}

func TestWebhookDispatcher_HandleEvent_ReturnsNilWhenListErrors(t *testing.T) {
	installRepo := &mockInstallationRepo{
		listErr: context.DeadlineExceeded, // simulate a DB error
	}
	versionRepo := newMockVersionRepo()
	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, nil)

	err := dispatcher.HandleEvent(context.Background(), makeEvent(t, "tenant-dberr", eventbus.ProductUpdated))
	if err != nil {
		t.Fatalf("HandleEvent should return nil even when list fails, got: %v", err)
	}
}

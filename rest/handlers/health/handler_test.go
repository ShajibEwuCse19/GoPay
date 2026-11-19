package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockPinger struct {
	pingFunc func(ctx context.Context) error
}

func (m *mockPinger) PingContext(ctx context.Context) error {
	if m.pingFunc != nil {
		return m.pingFunc(ctx)
	}
	return nil
}

func TestLive_AlwaysReturns200(t *testing.T) {
	h := NewHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	rec := httptest.NewRecorder()

	h.Live(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "UP" {
		t.Errorf("expected status 'UP', got %q", resp["status"])
	}
}

func TestReady_DatabaseUp_Returns200(t *testing.T) {
	pinger := &mockPinger{
		pingFunc: func(ctx context.Context) error {
			return nil
		},
	}
	h := NewHandler(pinger)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rec := httptest.NewRecorder()

	h.Ready(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp struct {
		Status string            `json:"status"`
		Checks map[string]string `json:"checks"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "UP" {
		t.Errorf("expected status 'UP', got %q", resp.Status)
	}
	if resp.Checks["database"] != "UP" {
		t.Errorf("expected checks.database 'UP', got %q", resp.Checks["database"])
	}
}

func TestReady_DatabaseDown_Returns503(t *testing.T) {
	pinger := &mockPinger{
		pingFunc: func(ctx context.Context) error {
			return errors.New("connection refused to postgres:5432")
		},
	}
	h := NewHandler(pinger)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rec := httptest.NewRecorder()

	h.Ready(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}

	var resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "DOWN" {
		t.Errorf("expected status 'DOWN', got %q", resp.Status)
	}
	if resp.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestReady_DatabaseTimeout_Returns503(t *testing.T) {
	pinger := &mockPinger{
		pingFunc: func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	}
	h := NewHandler(pinger)

	// Provide a short request context to simulate immediate timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.Ready(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}
}

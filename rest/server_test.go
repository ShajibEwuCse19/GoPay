package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ShajibEwuCse19/gopay-core/config"
	"github.com/ShajibEwuCse19/gopay-core/domain"
	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/customer"
	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/health"
)

type dummyService struct{}

func (d *dummyService) Create(ctx context.Context, in domain.NewCustomer) (*domain.Customer, error) {
	return nil, nil
}

type dummyPinger struct {
	err error
}

func (d *dummyPinger) PingContext(ctx context.Context) error {
	return d.err
}

func TestServerPipeline(t *testing.T) {
	cfg := &config.Config{
		HTTP: config.HTTPConfig{
			Port: 8080,
		},
		Timeouts: config.TimeoutConfig{
			CreateCustomer: time.Second,
		},
	}

	h := customer.NewHandler(cfg, &dummyService{})
	healthH := health.NewHandler(&dummyPinger{})
	server := NewServer(h, healthH, cfg)
	handler := server.buildHandler()

	t.Run("returns X-Request-ID header and request_id in error body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/customers/register", strings.NewReader(`{"invalid":true}`))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		headerID := rec.Header().Get("X-Request-ID")
		if headerID == "" {
			t.Fatal("expected X-Request-ID header in response, got empty")
		}

		var resp map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to parse json response: %v", err)
		}

		if resp["request_id"] != headerID {
			t.Errorf("body request_id %q does not match header %q", resp["request_id"], headerID)
		}
	})

	t.Run("serves /health/live with 200 OK", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 for live probe, got %d", rec.Code)
		}
	})

	t.Run("serves /health/ready with 200 OK when healthy", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 for ready probe, got %d", rec.Code)
		}
	})

	t.Run("serves /health/ready with 503 when DB down", func(t *testing.T) {
		serverDown := NewServer(h, health.NewHandler(&dummyPinger{err: errors.New("db down")}), cfg)
		handlerDown := serverDown.buildHandler()

		req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
		rec := httptest.NewRecorder()

		handlerDown.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("expected 503 for ready probe when db down, got %d", rec.Code)
		}
	})
}

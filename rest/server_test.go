package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ShajibEwuCse19/gopay-core/config"
	"github.com/ShajibEwuCse19/gopay-core/domain"
	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/customer"
)

type dummyService struct{}

func (d *dummyService) Create(ctx context.Context, in domain.NewCustomer) (*domain.Customer, error) {
	return nil, nil
}

func TestServerPipeline_RequestID(t *testing.T) {
	cfg := &config.Config{
		HTTP: config.HTTPConfig{
			Port: 8080,
		},
		Timeouts: config.TimeoutConfig{
			CreateCustomer: time.Second,
		},
	}

	h := customer.NewHandler(cfg, &dummyService{})
	server := NewServer(h, cfg)
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
}

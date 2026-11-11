package customer

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ShajibEwuCse19/gopay-core/config"
	"github.com/ShajibEwuCse19/gopay-core/domain"
)

type fakeService struct {
	result *domain.Customer
	err    error
}

func (f *fakeService) Create(ctx context.Context, in domain.NewCustomer) (*domain.Customer, error) {
	return f.result, f.err
}

func newTestHandler(svc Service) *Handler {
	cfg := &config.Config{Timeouts: config.TimeoutConfig{CreateCustomer: time.Second}}
	return NewHandler(cfg, svc)
}

const validBody = `{"phone_number":"+8801712345678","email":"rahim@example.com","full_name":"Rahim","pin":"482915"}`

func post(h *Handler, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/customers/register", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.CreateCustomers(rec, req)
	return rec
}

func TestCreateCustomers_Success(t *testing.T) {
	h := newTestHandler(&fakeService{result: &domain.Customer{
		ID:          "id-1",
		PhoneNumber: "+8801712345678",
		Email:       "rahim@example.com",
		FullName:    "Rahim",
		PinHash:     "secret-hash",
		KycStatus:   "UNVERIFIED",
		Status:      "ACTIVE",
	}})

	rec := post(h, validBody)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	body := rec.Body.String()
	if strings.Contains(body, "secret-hash") || strings.Contains(body, "pin") {
		t.Errorf("response leaks the PIN or its hash: %s", body)
	}
	if !strings.Contains(body, `"id":"id-1"`) {
		t.Errorf("response is missing the customer id: %s", body)
	}
}

func TestCreateCustomers_Errors(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcErr     error
		wantStatus int
	}{
		{"invalid json", `{not json`, nil, http.StatusBadRequest},
		{"unknown field", `{"phone_number":"+8801712345678","pin_hash":"x"}`, nil, http.StatusBadRequest},
		{"body too large", `{"full_name":"` + strings.Repeat("a", maxRequestBodyBytes) + `"}`, nil, http.StatusRequestEntityTooLarge},
		{"validation error", validBody, &domain.ValidationError{Field: "pin", Message: "is too easy to guess"}, http.StatusBadRequest},
		{"duplicate customer", validBody, domain.ErrCustomerAlreadyExists, http.StatusConflict},
		{"timeout", validBody, context.DeadlineExceeded, http.StatusGatewayTimeout},
		{"wrapped timeout", validBody, errors.Join(errors.New("db"), context.DeadlineExceeded), http.StatusGatewayTimeout},
		{"unknown error", validBody, errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := post(newTestHandler(&fakeService{err: tt.svcErr}), tt.body)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), `"error"`) {
				t.Errorf("error body is not in the JSON error shape: %s", rec.Body.String())
			}
		})
	}
}

func TestCreateCustomers_ClientCancelWritesNoBody(t *testing.T) {
	rec := post(newTestHandler(&fakeService{err: context.Canceled}), validBody)

	if rec.Body.Len() != 0 {
		t.Errorf("want empty body when the client is gone, got %s", rec.Body.String())
	}
}

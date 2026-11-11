package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ShajibEwuCse19/gopay-core/domain"
	"golang.org/x/crypto/bcrypt"
)

type fakeRepo struct {
	called   bool
	received domain.Customer
	err      error
}

func (f *fakeRepo) Create(ctx context.Context, c domain.Customer) (*domain.Customer, error) {
	f.called = true
	f.received = c
	if f.err != nil {
		return nil, f.err
	}
	return &c, nil
}

func validInput() domain.NewCustomer {
	return domain.NewCustomer{
		PhoneNumber: "+8801712345678",
		Email:       "rahim@example.com",
		FullName:    "Rahim Uddin",
		Pin:         "482915",
	}
}

func TestCreate_HashesPinAndNormalizes(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo)

	in := validInput()
	in.Email = "  Rahim@Example.COM "
	in.FullName = "  Rahim Uddin  "

	if _, err := svc.Create(context.Background(), in); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.received.PinHash == in.Pin {
		t.Fatal("PIN was stored in plain text")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.received.PinHash), []byte(in.Pin)); err != nil {
		t.Fatalf("stored value is not a valid bcrypt hash of the PIN: %v", err)
	}
	if repo.received.Email != "rahim@example.com" {
		t.Errorf("email not normalized: %q", repo.received.Email)
	}
	if repo.received.FullName != "Rahim Uddin" {
		t.Errorf("full name not trimmed: %q", repo.received.FullName)
	}
}

func TestCreate_Validation(t *testing.T) {
	tests := []struct {
		name  string
		edit  func(*domain.NewCustomer)
		field string
	}{
		{"phone without plus", func(c *domain.NewCustomer) { c.PhoneNumber = "8801712345678" }, "phone_number"},
		{"phone too short", func(c *domain.NewCustomer) { c.PhoneNumber = "+123" }, "phone_number"},
		{"email empty", func(c *domain.NewCustomer) { c.Email = "" }, "email"},
		{"email invalid", func(c *domain.NewCustomer) { c.Email = "not-an-email" }, "email"},
		{"email with display name", func(c *domain.NewCustomer) { c.Email = "Rahim <rahim@example.com>" }, "email"},
		{"name empty", func(c *domain.NewCustomer) { c.FullName = "   " }, "full_name"},
		{"pin too short", func(c *domain.NewCustomer) { c.Pin = "48291" }, "pin"},
		{"pin has letters", func(c *domain.NewCustomer) { c.Pin = "48291a" }, "pin"},
		{"pin repeated", func(c *domain.NewCustomer) { c.Pin = "000000" }, "pin"},
		{"pin ascending", func(c *domain.NewCustomer) { c.Pin = "123456" }, "pin"},
		{"pin descending", func(c *domain.NewCustomer) { c.Pin = "654321" }, "pin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			in := validInput()
			tt.edit(&in)

			_, err := NewService(repo).Create(context.Background(), in)

			var vErr *domain.ValidationError
			if !errors.As(err, &vErr) {
				t.Fatalf("want ValidationError, got %v", err)
			}
			if vErr.Field != tt.field {
				t.Errorf("field = %q, want %q", vErr.Field, tt.field)
			}
			if repo.called {
				t.Error("repo must not be called when input is invalid")
			}
		})
	}
}

func TestCreate_CancelledContextSkipsWork(t *testing.T) {
	repo := &fakeRepo{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewService(repo).Create(ctx, validInput())

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}
	if repo.called {
		t.Error("repo must not be called with a cancelled context")
	}
}

func TestCreate_PassesRepoErrorUp(t *testing.T) {
	repo := &fakeRepo{err: domain.ErrCustomerAlreadyExists}

	_, err := NewService(repo).Create(context.Background(), validInput())

	if !errors.Is(err, domain.ErrCustomerAlreadyExists) {
		t.Fatalf("want ErrCustomerAlreadyExists, got %v", err)
	}
}

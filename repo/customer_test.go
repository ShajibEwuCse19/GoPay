package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/ShajibEwuCse19/gopay-core/domain"
	"github.com/lib/pq"
)

func TestMapCreateError(t *testing.T) {
	other := errors.New("connection reset")

	tests := []struct {
		name string
		in   error
		want error
	}{
		{"unique violation", &pq.Error{Code: "23505"}, domain.ErrCustomerAlreadyExists},
		{"wrapped unique violation", errors.Join(errors.New("insert"), &pq.Error{Code: "23505"}), domain.ErrCustomerAlreadyExists},
		{"other postgres error stays as is", &pq.Error{Code: "23502"}, nil},
		{"context deadline stays as is", context.DeadlineExceeded, context.DeadlineExceeded},
		{"unknown error stays as is", other, other},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapCreateError(tt.in)

			if tt.want == nil {
				if errors.Is(got, domain.ErrCustomerAlreadyExists) {
					t.Fatalf("must not map to ErrCustomerAlreadyExists, got %v", got)
				}
				return
			}
			if !errors.Is(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

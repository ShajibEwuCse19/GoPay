package customer

import (
	"context"

	"github.com/ShajibEwuCse19/gopay-core/domain"
)

type Service interface {
	Create(ctx context.Context, c domain.Customer) (*domain.Customer, error)
}
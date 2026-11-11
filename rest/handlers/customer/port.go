package customer

import (
	"context"

	"github.com/ShajibEwuCse19/gopay-core/domain"
)

type Service interface {
	Create(ctx context.Context, in domain.NewCustomer) (*domain.Customer, error)
}

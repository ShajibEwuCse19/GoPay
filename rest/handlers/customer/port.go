package customer

import (
	"context"

	"github.com/ShajibEwuCse19/gopay-core/domain"
)

type Service interface {
	// input domain.NewCustomer and return domain.Customer
	// data security as pinHash will not be returned and domain.Customer is the main DB representation
	// We are not storing plain pin in the database model
	Create(ctx context.Context, in domain.NewCustomer) (*domain.Customer, error)
}

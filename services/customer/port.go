package services

import (
	"context"

	"github.com/ShajibEwuCse19/gopay-core/domain"
	"github.com/ShajibEwuCse19/gopay-core/rest/handlers/customer"
)

type Service interface {
	customer.Service
}

type CustomerRepo interface {
	Create(ctx context.Context, c domain.Customer) (*domain.Customer, error)
}

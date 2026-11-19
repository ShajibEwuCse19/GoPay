package services

import (
	"context"

	"github.com/ShajibEwuCse19/gopay-core/domain"
)

type CustomerRepo interface {
	Create(ctx context.Context, c domain.Customer) (*domain.Customer, error)
}

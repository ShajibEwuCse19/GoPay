package services

import (
	"context"

	"github.com/ShajibEwuCse19/gopay-core/domain"
)

type service struct {
	customerRepo CustomerRepo
}

func NewService(customerRepo CustomerRepo) Service {
	return &service{
		customerRepo: customerRepo,
	}
}

func (s *service) Create(ctx context.Context, c domain.Customer) (*domain.Customer, error) {
	customer, err := s.customerRepo.Create(ctx, c)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

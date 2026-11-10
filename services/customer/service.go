package services

import "github.com/ShajibEwuCse19/gopay-core/domain"

type service struct {
	customerRepo CustomerRepo
}

func NewService(customerRepo CustomerRepo) Service {
	return &service{
		customerRepo: customerRepo,
	}
}

func (s *service) Create(c domain.Customer) (*domain.Customer, error) {
	customer, err := s.customerRepo.Create(c)
	if err != nil {
		return nil, err
	}
	return customer, nil
}


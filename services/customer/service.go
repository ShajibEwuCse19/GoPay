package services

import (
	"context"

	"github.com/ShajibEwuCse19/gopay-core/domain"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	customerRepo CustomerRepo
}

func NewService(customerRepo CustomerRepo) *Service {
	return &Service{
		customerRepo: customerRepo,
	}
}

func (s *Service) Create(ctx context.Context, in domain.NewCustomer) (*domain.Customer, error) {
	in = normalize(in)
	if err := validate(in); err != nil {
		return nil, err
	}

	// bcrypt is slow on purpose; don't start it if the request is already over
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	pinHash, err := bcrypt.GenerateFromPassword([]byte(in.Pin), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var emailPtr *string
	if in.Email != "" {
		emailPtr = &in.Email
	}

	return s.customerRepo.Create(ctx, domain.Customer{
		PhoneNumber: in.PhoneNumber,
		Email:       emailPtr,
		FullName:    in.FullName,
		PinHash:     string(pinHash),
	})
}

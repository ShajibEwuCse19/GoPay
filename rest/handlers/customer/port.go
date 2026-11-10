package customer

import "github.com/ShajibEwuCse19/gopay-core/domain"

type Service interface {
	Create(c domain.Customer) (*domain.Customer, error)
}
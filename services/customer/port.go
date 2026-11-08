package services

import "github.com/ShajibEwuCse19/gopay-core/rest/handlers/customer"

type Service interface {
	customer.Service
}

type CustomerRepo interface {
}
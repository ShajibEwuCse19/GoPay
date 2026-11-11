package domain

import (
	"errors"
	"time"
)

var ErrCustomerAlreadyExists = errors.New("customer already exists")

// ValidationError is returned when user input breaks a business rule (maps to HTTP 400).
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// NewCustomer is the input to register a customer. Pin is the plain PIN; it is hashed in the service.
type NewCustomer struct {
	PhoneNumber string
	Email       string
	FullName    string
	Pin         string
}

type Customer struct {
	ID          string    `json:"id" db:"id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Email       string    `json:"email" db:"email"`
	FullName    string    `json:"full_name" db:"full_name"`
	PinHash     string    `json:"-" db:"pin_hash"` // Hide this field from JSON output
	KycStatus   string    `json:"kyc_status" db:"kyc_status"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

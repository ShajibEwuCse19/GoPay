package domain

import "time"

type Customer struct {
	ID          string    `json:"id" db:"id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Email       string    `json:"email" db:"email"`
	FullName    string    `json:"full_name" db:"full_name"`
	PinHash     string    `json:"pin_hash" db:"pin_hash"`
	KycStatus   string    `json:"kyc_status" db:"kyc_status"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
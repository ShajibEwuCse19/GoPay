package services

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/ShajibEwuCse19/gopay-core/domain"
)

const (
	pinLength      = 6
	maxEmailLength = 255
	maxFullNameLen = 120
)

// same rule as the chk_customers_phone_e164 constraint in the migration
var phoneE164 = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)

// normalize trims spaces and lowercases the email so stored data is tidy.
func normalize(c domain.NewCustomer) domain.NewCustomer {
	return domain.NewCustomer{
		PhoneNumber: strings.TrimSpace(c.PhoneNumber),
		Email:       strings.ToLower(strings.TrimSpace(c.Email)),
		FullName:    strings.TrimSpace(c.FullName),
		Pin:         c.Pin, // never trim or change a secret
	}
}

func validate(c domain.NewCustomer) error {
	if !phoneE164.MatchString(c.PhoneNumber) {
		return &domain.ValidationError{Field: "phone_number", Message: "must be in international format, e.g. +8801712345678"}
	}

	if c.Email == "" || len(c.Email) > maxEmailLength {
		return &domain.ValidationError{Field: "email", Message: "is required and must be at most 255 characters"}
	}
	addr, err := mail.ParseAddress(c.Email)
	if err != nil || addr.Address != c.Email {
		return &domain.ValidationError{Field: "email", Message: "is not a valid email address"}
	}

	if n := utf8.RuneCountInString(c.FullName); n == 0 || n > maxFullNameLen {
		return &domain.ValidationError{Field: "full_name", Message: "is required and must be at most 120 characters"}
	}

	return validatePin(c.Pin)
}

func validatePin(pin string) error {
	if len(pin) != pinLength {
		return &domain.ValidationError{Field: "pin", Message: "must be exactly 6 digits"}
	}
	for _, r := range pin {
		if r < '0' || r > '9' {
			return &domain.ValidationError{Field: "pin", Message: "must be exactly 6 digits"}
		}
	}
	if isRepeated(pin) || isSequential(pin) {
		return &domain.ValidationError{Field: "pin", Message: "is too easy to guess"}
	}
	return nil
}

// isRepeated reports PINs like 000000 or 777777.
func isRepeated(pin string) bool {
	return strings.Count(pin, string(pin[0])) == len(pin)
}

// isSequential reports PINs like 123456 or 654321.
func isSequential(pin string) bool {
	asc, desc := true, true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[i-1]+1 {
			asc = false
		}
		if pin[i] != pin[i-1]-1 {
			desc = false
		}
	}
	return asc || desc
}

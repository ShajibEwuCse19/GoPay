package domain

import "errors"

// Sentinel Domain Errors needed for current features
var (
	ErrConflict              = errors.New("resource already exists")
	ErrCustomerAlreadyExists = ErrConflict // Alias for backward compatibility
)

// ValidationError handles field-level invalid inputs (maps to HTTP 400).
type ValidationError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

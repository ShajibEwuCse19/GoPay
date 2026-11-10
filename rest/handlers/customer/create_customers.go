package customer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ShajibEwuCse19/gopay-core/domain"
	"github.com/ShajibEwuCse19/gopay-core/utils"
)

const createCustomerTimeout = 5 * time.Second

type ReqCreateCustomer struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PinHash     string `json:"pin_hash"`
}

func (h *Handler) CreateCustomers(w http.ResponseWriter, r *http.Request) {
	var req ReqCreateCustomer
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request body. Give me valid json", http.StatusBadRequest)
		return
	}

	if req.PhoneNumber == "" || req.Email == "" || req.FullName == "" || req.PinHash == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// r.Context() is cancelled if the client disconnects; the timeout adds a time limit on top.
	ctx, cancel := context.WithTimeout(r.Context(), createCustomerTimeout)
	defer cancel()

	createCustomer, err := h.svc.Create(ctx, domain.Customer{
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		FullName:    req.FullName,
		PinHash:     req.PinHash,
	})

	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, "Request timed out. Please try again", http.StatusGatewayTimeout)
		case errors.Is(err, context.Canceled):
			fmt.Println("register: client cancelled the request") // client is gone, nobody to answer
		case errors.Is(err, domain.ErrCustomerAlreadyExists):
			http.Error(w, "Phone number or email already registered", http.StatusConflict)
		default:
			http.Error(w, "Error creating customer", http.StatusInternalServerError)
		}
		return
	}

	utils.SendResponse(w, http.StatusOK, createCustomer)
}
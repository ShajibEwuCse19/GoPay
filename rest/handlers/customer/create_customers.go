package customer

import (
	"encoding/json"
	"net/http"

	"github.com/ShajibEwuCse19/gopay-core/domain"
	"github.com/ShajibEwuCse19/gopay-core/utils"
)

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

	createCustomer, err := h.svc.Create(domain.Customer{
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		FullName:    req.FullName,
		PinHash:     req.PinHash,
	})

	if err != nil {
		http.Error(w, "Error creating customer", http.StatusInternalServerError)
		return
	}

	utils.SendResponse(w, http.StatusOK, createCustomer)
}
package customer

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ShajibEwuCse19/gopay-core/domain"
	"github.com/ShajibEwuCse19/gopay-core/utils"
)

// maxRequestBodyBytes limits how much of the request body we read (1 MB).
const maxRequestBodyBytes = 1 << 20

type ReqCreateCustomer struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Pin         string `json:"pin"`
}

type RespCustomer struct {
	ID          string    `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	KycStatus   string    `json:"kyc_status"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (h *Handler) CreateCustomers(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var req ReqCreateCustomer
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			utils.SendError(w, http.StatusRequestEntityTooLarge, "Request body is too large")
			return
		}
		utils.SendError(w, http.StatusBadRequest, "Invalid request body. Give me valid json")
		return
	}

	// r.Context() is cancelled if the client disconnects; the timeout adds a time limit on top.
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeouts.CreateCustomer)
	defer cancel()

	created, err := h.svc.Create(ctx, domain.NewCustomer{
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		FullName:    req.FullName,
		Pin:         req.Pin,
	})
	if err != nil {
		utils.WriteError(w, r, err)
		return
	}

	utils.SendResponse(w, http.StatusCreated, RespCustomer{
		ID:          created.ID,
		PhoneNumber: created.PhoneNumber,
		Email:       created.Email,
		FullName:    created.FullName,
		KycStatus:   created.KycStatus,
		Status:      created.Status,
		CreatedAt:   created.CreatedAt,
	})
}

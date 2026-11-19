package customer

import (
	"context"
	"encoding/json"
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
		utils.WriteError(w, r, err)
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

	var emailStr string
	if created.Email != nil {
		emailStr = *created.Email
	}

	utils.SendResponse(w, http.StatusCreated, RespCustomer{
		ID:          created.ID,
		PhoneNumber: created.PhoneNumber,
		Email:       emailStr,
		FullName:    created.FullName,
		KycStatus:   created.KycStatus,
		Status:      created.Status,
		CreatedAt:   created.CreatedAt,
	})
}

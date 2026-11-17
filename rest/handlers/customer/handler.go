package customer

import "github.com/ShajibEwuCse19/gopay-core/config"

type Handler struct {
	cfg *config.Config
	svc Service
}

func NewHandler(cfg *config.Config, s Service) *Handler {
	return &Handler{
		cfg: cfg,
		svc: s,
	}
}

package handler

import (
	"net/http"
)

type Handler struct{}

func (h *Handler) InitRoutes() {
	http.HandleFunc("/api/orders", h.GetOrders)
}

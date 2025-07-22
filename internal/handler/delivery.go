package handler

import "net/http"

func (h *Handler) GetDeliveries(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /api/deliveries works"))
}

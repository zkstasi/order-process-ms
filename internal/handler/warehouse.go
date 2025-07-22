package handler

import "net/http"

func (h *Handler) GetWarehouses(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /api/warehouses works"))
}

package handler

import "net/http"

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /api/users works"))
}

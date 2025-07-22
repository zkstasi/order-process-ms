package handler

import "net/http"

// метод-обработчик для http-запроса GET

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /api/orders works")) // отправляет в тело http-ответа простую строку
}

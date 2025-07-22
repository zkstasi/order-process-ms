package handler

import (
	"net/http"
	"os"
)

// метод-обработчик для http-запроса GET

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("data/orders.json")
	if err != nil {
		http.Error(w, "Ошибка при чтении файла orders.json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

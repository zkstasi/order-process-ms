package handler

import "net/http"

// метод регистрирует маршруты и связывает их с обработчиками

func (h *Handler) InitRoutes() {
	h.mux.HandleFunc("/api/orders", h.GetOrders)
	h.mux.HandleFunc("/api/users", h.GetUsers)
	h.mux.HandleFunc("/api/deliveries", h.GetDeliveries)
	h.mux.HandleFunc("/api/warehouses", h.GetWarehouses)
}

// серверу должен передаваться именно такой интерфейс для обработки http-запросов
// возвращает объект, реализующий интерфейс, это нужно, чтобы передать его в ListenAndServe

func (h *Handler) Router() http.Handler {
	return h.mux
}

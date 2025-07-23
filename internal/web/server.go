package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"order-ms/internal/model"
	"order-ms/internal/repository"
)

type Server struct {
	address    string       // Адрес, по которому будет слушать сервер
	httpServer *http.Server // Указатель на стандартный http-сервер
}

// создание нового сервера

func NewServer(address string) *Server {
	return &Server{
		address: address,
		httpServer: &http.Server{
			Addr: address,
		},
	}
}

// метод запуска http-сервера

func (s *Server) Start() error {
	http.HandleFunc("/api/orders", s.handleOrders) // связь url с методом-обработчиком handleCreateOrder
	http.HandleFunc("/api/orders/", s.handleGetOrderByID)

	log.Printf("Server starting on %s\n", s.address)
	return s.httpServer.ListenAndServe() // запускает сервер и блокирует при ошибке
}

// Структура для парсинга, какие поля ожидаем в json-запросе
type createOrderRequest struct {
	UserID string `json:"user_id"`
}

// Обработчик POST-запроса ("ручка"), вызывается когда приходит запрос

func (s *Server) handleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// Обработчик POST-запроса ("ручка"), вызывается когда приходит запрос
	case "POST":
		body, err := io.ReadAll(r.Body) // r.Body - поток с данными от клиента
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		// распарсим user_id в структуру
		var req createOrderRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		// проверяем, что UserID не пустой (валидация)
		if req.UserID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}
		// создаем заказ и сохраняем
		order := model.NewOrder(req.UserID)
		repository.SaveStorable(order)
		// возвращаем результат клиенту
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)

	// обработчик GET-запроса
	case "GET":
		// получаем список всех заказов
		orders := repository.GetOrders()
		w.Header().Set("Content-Type", "application/json")
		// отправляем клиенту json-массив заказов
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			http.Error(w, "Failed to encode orders", http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Обработчик GET c ID
func (s *Server) handleGetOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// получаем id из пути /api/orders/{id}
	id := r.URL.Path[len("/api/orders/"):] // вырезаем часть после "/api/orders/"
	if id == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
	// ищем заказ
	order := repository.GetOrderByID(id)
	if order == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	// отправляем json
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode order", http.StatusInternalServerError)
		return
	}
}

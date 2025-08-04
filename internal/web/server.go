package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"order-ms/internal/model"
	"order-ms/internal/repository"
	"time"
)

type Server struct {
	address    string       // Адрес, по которому будет слушать сервер
	httpServer *http.Server // Указатель на стандартный http-сервер
}

// Структура для парсинга, какие поля ожидаем в json-запросе
type createOrderRequest struct {
	UserID string `json:"user_id"`
}

type createUserRequest struct {
	Name string `json:"name"`
}

// cтруктура для PUT-запроса (обновления)
type updateOrderStatus struct {
	Status model.OrderStatus `json:"status"`
}

type updateUserRequest struct {
	Name string `json:"name"`
}

// создание нового сервера

func NewServer(address string) *Server {
	mux := http.NewServeMux() // создаем локальный маршрутизатор

	s := &Server{
		address: address,
		httpServer: &http.Server{
			Addr:         address,
			Handler:      mux,
			ReadTimeout:  10 * time.Second, // сколько времени сервер ждёт запрос от клиента (например, тело запроса)
			WriteTimeout: 10 * time.Second, // сколько времени дается серверу на отправку ответа клиенту
			IdleTimeout:  60 * time.Second, // время ожидания между запросами, если клиент держит соединение открытым
		},
	}
	//регистрируем эндпоинты (маршруты) в mux, по которым будут обрабатываться запросы
	mux.HandleFunc("/api/orders", s.handleOrders) // связь url с методом-обработчиком
	mux.HandleFunc("/api/orders/", s.handleOrderByID)
	mux.HandleFunc("/api/orders/confirm/", s.handleOrderConfirm)
	mux.HandleFunc("/api/orders/delivery/", s.handleOrderDelivery)
	mux.HandleFunc("/api/orders/cancel/", s.handleOrderCancel)
	mux.HandleFunc("/api/users", s.handleUsers)
	mux.HandleFunc("/api/users/", s.handleUserByID)

	return s
}

// метод запуска http-сервера

func (s *Server) Start() error {
	log.Printf("Server starting on %s\n", s.address)
	return s.httpServer.ListenAndServe() // запускает сервер и блокирует при ошибке
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

func (s *Server) handleOrderByID(w http.ResponseWriter, r *http.Request) {
	// получаем id из пути /api/orders/{id}
	id := r.URL.Path[len("/api/orders/"):] // вырезаем часть после "/api/orders/"
	if id == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
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
	case "PUT":
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		var req updateOrderStatus
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		ok := repository.ConfirmOrder(id)
		if !ok {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case "DELETE":
		ok := repository.DeleteOrder(id)
		if !ok {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Обработчик запрос POST на подтверждение заказа
func (s *Server) handleOrderConfirm(w http.ResponseWriter, r *http.Request) {
	// извлекаем id заказа
	id := r.URL.Path[len("/api/orders/confirm/"):]
	//проверка: передан ли id
	if id == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "POST":
		// находим заказ в хранилище по его id
		order := repository.GetOrderByID(id)
		// если заказ не найден, возвращаем ошибку
		if order == nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		// проверяем, можно ли подтвердить заказ
		if order.Status != model.OrderCreated {
			http.Error(w, "Order must be in 'created' status to confirm", http.StatusBadRequest)
			return
		}
		// меняем статус заказа на "1"
		repository.ConfirmOrder(id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleOrderDelivery(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/orders/delivery/"):]
	if id == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "POST":
		order := repository.GetOrderByID(id)
		if order == nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		if order.Status != model.OrderConfirmed {
			http.Error(w, "Order must be in 'confirmed' status to be delivered", http.StatusBadRequest)
			return
		}
		repository.DeliveredOrder(id)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleOrderCancel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/orders/cancel/"):]
	if id == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "POST":
		order := repository.GetOrderByID(id)
		if order == nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		if order.Status != model.OrderCreated && order.Status != model.OrderConfirmed {
			http.Error(w, "Only orders with status 'created' or 'confirmed' can be cancelled", http.StatusBadRequest)
			return
		}
		ok := repository.CancelOrder(id)
		if !ok {
			http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		var req createUserRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}
		user := model.NewUser(req.Name)
		repository.SaveStorable(user)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Failed to encode user", http.StatusInternalServerError)
			return
		}
	case "GET":
		users := repository.GetUsers()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Failed to encode users", http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func (s *Server) handleUserByID(w http.ResponseWriter, r *http.Request) {
	// Извлекаем id из URL (/api/users/{id})
	id := r.URL.Path[len("/api/users/"):]
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		// ищем пользователя по id
		user := repository.GetUserByID(id)
		if user == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		// кодируем пользователя в json и отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Failed to encode user", http.StatusInternalServerError)
			return
		}
	case "PUT":
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		var req updateUserRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		ok := repository.UpdateUserName(id, req.Name)
		if !ok {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case "DELETE":
		ok := repository.DeleteUser(id)
		if !ok {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

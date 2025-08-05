package web

import (
	"encoding/json"
	httpSwagger "github.com/swaggo/http-swagger"
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
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return s
}

// метод запуска http-сервера

func (s *Server) Start() error {
	log.Printf("Server starting on %s\n", s.address)
	return s.httpServer.ListenAndServe() // запускает сервер и блокирует при ошибке
}

// handleOrders обрабатывает запросы к заказам
// @Summary Создать заказ
// @Description Создает новый заказ пользователя с переданным userID
// @Tags Orders
// @Accept json
// @Produce json
// @Param user body createOrderRequest true "User ID"
// @Success 201 {object} model.Order "Созданный заказ"
// @Failure 400 {object} object "Неверный JSON или не указан user ID"
// @Failure 500 {object} object "Ошибка кодирования ответа"
// @Router /api/orders [post]

// @Summary Получить список заказов
// @Description Возвращает все созданные заказы
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200 {array} model.Order "Список заказов"
// @Failure 500 {object} object "Ошибка кодирования ответа"
// @Router /api/orders [get]
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

// handleOrderByID обрабатывает запросы к заказу по его ID
// @Summary Получить заказ по ID
// @Description Возвращает заказ с указанным ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "ID заказа"
// @Success 200 {object} model.Order "Найденный заказ"
// @Failure 400 {object} object "Некорректный ID"
// @Failure 404 {object} object "Заказ не найден"
// @Router /api/orders/{id} [get]

// @Summary Удалить заказ
// @Description Удаляет заказ с указанным ID
// @Tags Orders
// @Param id path string true "ID заказа"
// @Success 204 "Заказ успешно удален"
// @Failure 400 {object} object "Некорректный ID"
// @Failure 404 {object} object "Заказ не найден"
// @Router /api/orders/{id} [delete]
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

// handleOrderConfirm подтверждает заказ складом и переводит его в статус "подтвержден" (1)
// @Summary Подтверждение заказа
// @Description Подтверждает заказ, если он находится в статусе "создан" (0)
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "ID заказа"
// @Success 204 "No Content - успешное подтверждение"
// @Failure 400 {object} object "Некорректный запрос или статус заказа не позволяет подтверждение"
// @Failure 404 {object} object "Заказ не найден"
// @Router /api/orders/confirm/{id} [post]
func (s *Server) handleOrderConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// извлекаем id заказа
	id := r.URL.Path[len("/api/orders/confirm/"):]
	//проверка: передан ли id
	if id == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
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
}

// handleOrderDelivery переводит заказ в статус "доставлен" (2)
// @Summary Отметить заказ как доставленный
// @Description Переводит заказ в статус "доставлен", если он в статусе "подтвержден"
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "ID заказа"
// @Success 204 "No Content - заказ успешно доставлен"
// @Failure 400 {object} object "Некорректный запрос или заказ в нужном статусе"
// @Failure 404 {object} object "Заказ не найден"
// @Failure 405 {object} object "Метод не поддерживается"
// @Router /api/orders/delivery/{id} [post]
func (s *Server) handleOrderDelivery(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Path[len("/api/orders/delivery/"):]
	if id == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
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
}

// handleOrderCancel отменяет заказ
// @Summary Отмена заказа
// @Description Отменяет заказ, если он в статусе "создан" или "подтвержден"
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "ID заказа"
// @Success 204 "No Content - заказ успешно отменен"
// @Failure 400 {object} object "Некорректный статус заказа для отмены"
// @Failure 404 {object} object "Заказ не найден"
// @Failure 405 {object} object "Метод не поддерживается"
// @Router /api/orders/cancel/{id} [post]
func (s *Server) handleOrderCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Path[len("/api/orders/cancel/"):]
	if id == "" {
		http.Error(w, "Missing order ID", http.StatusBadRequest)
		return
	}
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
}

// handleUsers обрабатывает запросы с пользователем
// @Summary Создать пользователя
// @Description Создает нового пользователя с переданным именем
// @Tags Users
// @Accept json
// @Produce json
// @Param user body createUserRequest true "Имя пользователя"
// @Success 200 {object} model.User "Созданный пользователь"
// @Failure 400 {object} object "Неверный JSON или не указано имя"
// @Failure 500 {object} object "Ошибка кодирования ответа"
// @Router /api/users [post]

// @Summary Получить список пользователей
// @Description Возвращает всех зарегистрированных пользователей
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} model.User "Список пользователей"
// @Failure 500 {object} object "Ошибка кодирования ответа"
// @Router /api/users [get]
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

// handleUserByID обрабатывает операции с пользователем по ID
// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по указанному ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} model.User "Пользователь найден"
// @Failure 400 {object} object "Отсутствует или некорректный ID"
// @Failure 404 {object} object "Пользователь не найден"
// @Failure 500 {object} object "Ошибка кодирования ответа"
// @Router /api/users/{id} [get]

// @Summary Обновить имя пользователя по ID
// @Description Обновляет имя пользователя по указанному ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Param user body updateUserRequest true "Новое имя пользователя"
// @Success 204 "Успешное обновление без тела ответа"
// @Failure 400 {object} object "Неверный JSON или некорректные данные"
// @Failure 404 {object} object "Пользователь не найден"
// @Router /api/users/{id} [put]

// @Summary Удалить пользователя по ID
// @Description Удаляет пользователя по указанному ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 204 "Пользователь успешно удалён"
// @Failure 404 {object} object "Пользователь не найден"
// @Router /api/users/{id} [delete]
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
		if req.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
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

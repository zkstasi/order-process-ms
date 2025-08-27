package web

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"order-ms/internal/model"
	"order-ms/internal/service"
	"time"
)

type Server struct {
	address    string       // Адрес, по которому будет слушать сервер
	httpServer *http.Server // Указатель на стандартный http-сервер
	repo       service.Repository
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

func NewServer(address string, repo service.Repository) *Server {
	router := gin.New()

	s := &Server{
		address: address,
		httpServer: &http.Server{
			Addr:         address,
			Handler:      router,
			ReadTimeout:  10 * time.Second, // сколько времени сервер ждёт запрос от клиента (например, тело запроса)
			WriteTimeout: 10 * time.Second, // сколько времени дается серверу на отправку ответа клиенту
			IdleTimeout:  60 * time.Second, // время ожидания между запросами, если клиент держит соединение открытым
		},
		repo: repo,
	}
	//регистрируем эндпоинты (маршруты) в gin, по которым будут обрабатываться запросы
	router.POST("/api/orders", s.handleOrderCreate) // связь url с методом-обработчиком
	router.GET("/api/orders", s.handleOrderList)
	router.GET("/api/orders/:id", s.handleOrderGetByID)
	router.DELETE("/api/orders/:id", s.handleOrderDeleteByID)
	router.POST("/api/orders/confirm/:id", s.handleOrderConfirm)
	router.POST("/api/orders/delivery/:id", s.handleOrderDelivery)
	router.POST("/api/orders/cancel/:id", s.handleOrderCancel)

	router.POST("/api/users", s.handleUserCreate)
	router.GET("/api/users", s.handleUserList)
	router.GET("/api/users/:id", s.handleUserGetByID)
	router.PUT("/api/users/:id", s.handleUserUpdateByID)
	router.DELETE("/api/users/:id", s.handleUserDeleteByID)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return s
}

// метод запуска http-сервера
func (s *Server) Start() error {
	log.Printf("Server starting on %s\n", s.address)
	return s.httpServer.ListenAndServe() // запускает сервер и блокирует при ошибке
}

// handleOrderCreate создает новый заказ
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
func (s *Server) handleOrderCreate(c *gin.Context) {
	// распарсим user_id в структуру
	var req createOrderRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	// проверяем, что UserID не пустой (валидация)
	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	// создаем заказ и сохраняем
	order := model.NewOrder(req.UserID)
	if err := s.repo.Save(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot save order"})
		return
	}
	// возвращаем результат клиенту
	c.JSON(http.StatusCreated, order)
}

// handleOrderList формирует список всех заказов
// @Summary Получить список заказов
// @Description Возвращает все созданные заказы
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200 {array} model.Order "Список заказов"
// @Failure 500 {object} object "Ошибка кодирования ответа"
// @Router /api/orders [get]
func (s *Server) handleOrderList(c *gin.Context) {
	// получаем список всех заказов
	orders, err := s.repo.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot get orders"})
		return
	}
	// отправляем клиенту json-массив заказов
	c.JSON(http.StatusOK, orders)
}

// handleOrderGetByID получает заказ по его ID
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
func (s *Server) handleOrderGetByID(c *gin.Context) {
	// получаем id из пути /api/orders/{id}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order ID"})
		return
	}
	// ищем заказ
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot get order"})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// handleOrderDeleteByID удаляет заказ по его ID
// @Summary Удалить заказ
// @Description Удаляет заказ с указанным ID
// @Tags Orders
// @Param id path string true "ID заказа"
// @Success 204 "Заказ успешно удален"
// @Failure 400 {object} object "Некорректный ID"
// @Failure 404 {object} object "Заказ не найден"
// @Router /api/orders/{id} [delete]
func (s *Server) handleOrderDeleteByID(c *gin.Context) {
	// получаем id из пути /api/orders/{id}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order ID"})
		return
	}
	ok, err := s.repo.DeleteOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot delete order"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.Status(http.StatusNoContent)
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
func (s *Server) handleOrderConfirm(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order ID"})
		return
	}

	// подтверждаем заказ через репозиторий
	ok, err := s.repo.ConfirmOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm order"})
		return
	}

	if !ok {
		c.JSON(http.StatusConflict, gin.H{"error": "Order not found or not in CREATED status"})
		return
	}

	// берём обновлённый заказ
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order"})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found after confirm"})
		return
	}

	c.JSON(http.StatusOK, order)
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
func (s *Server) handleOrderDelivery(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order ID"})
		return
	}
	// помечаем заказ как доставленный
	ok, err := s.repo.DeliverOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark order as delivered"})
		return
	}
	if !ok {
		c.JSON(http.StatusConflict, gin.H{"error": "Order not found or not in CONFIRMED status"})
		return
	}

	// достаём обновлённый заказ
	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated order"})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found after delivery"})
		return
	}

	c.JSON(http.StatusOK, order)
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
func (s *Server) handleOrderCancel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order ID"})
		return
	}
	ok, err := s.repo.CancelOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}
	if !ok {
		c.JSON(http.StatusConflict, gin.H{"error": "Order not found or cannot be canceled"})
		return
	}

	// успешная отмена → возвращаем 204 No Content
	c.Status(http.StatusNoContent)
}

// handleUserCreate создает нового пользователя
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
func (s *Server) handleUserCreate(c *gin.Context) {
	var req createUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	// создаем пользователя и сохраняем
	user := model.NewUser(req.Name)
	if err := s.repo.Save(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot save user"})
		return
	}
	// возвращаем результат клиенту
	c.JSON(http.StatusCreated, user)
}

// handleUserList формирует список всех пользователей
// @Summary Получить список пользователей
// @Description Возвращает всех зарегистрированных пользователей
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} model.User "Список пользователей"
// @Failure 500 {object} object "Ошибка кодирования ответа"
// @Router /api/users [get]
func (s *Server) handleUserList(c *gin.Context) {
	users, err := s.repo.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot get users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// handleUserGetByID ищет пользователя по ID
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
func (s *Server) handleUserGetByID(c *gin.Context) {
	// Извлекаем id из URL (/api/users/{id})
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}
	// ищем пользователя по id
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// handleUserUpdateByID обновляет имя пользователя по ID
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
func (s *Server) handleUserUpdateByID(c *gin.Context) {
	// Извлекаем id из URL (/api/users/{id})
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}
	var req updateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	ok, err := s.repo.UpdateUserName(id, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	updatedUser, err := s.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// handleUserDeleteByID удаляет пользователя по ID
// @Summary Удалить пользователя по ID
// @Description Удаляет пользователя по указанному ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 204 "Пользователь успешно удалён"
// @Failure 404 {object} object "Пользователь не найден"
// @Router /api/users/{id} [delete]
func (s *Server) handleUserDeleteByID(c *gin.Context) {
	// Извлекаем id из URL (/api/users/{id})
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}
	ok, err := s.repo.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

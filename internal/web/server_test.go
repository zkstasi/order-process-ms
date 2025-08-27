package web

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"order-ms/internal/model"
	"order-ms/internal/repository"
	"strings"
	"testing"
)

// тест ручки GET для получения заказов
func TestGetOrders(t *testing.T) {
	gin.SetMode(gin.TestMode) // чтобы не было лишних логов

	// загружаем данные из файлов перед тестом
	if err := repository.LoadOrdersFromFile("../../data/orders.json"); err != nil {
		t.Fatalf("Не удалось загрузить заказы: %v", err)
	}
	if err := repository.LoadUsersFromFile("../../data/users.json"); err != nil {
		t.Fatalf("Не удалось загрузить пользователей: %v", err)
	}

	tests := []struct {
		name        string
		expectedMin int
	}{
		{
			name:        "orders exist",
			expectedMin: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// создаем сервер
			s := NewServer(":8080")

			// получаем роутер
			r := s.httpServer.Handler.(*gin.Engine)

			// создаем запрос
			req, _ := http.NewRequest("GET", "/api/orders", nil)
			w := httptest.NewRecorder()

			// выполняем запрос
			r.ServeHTTP(w, req)

			// проверяем код ответа
			assert.Equal(t, http.StatusOK, w.Code)

			// проверяем тело ответа
			var got []model.Order
			err := json.Unmarshal(w.Body.Bytes(), &got)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(got), tc.expectedMin)
		})
	}
}

// тест ручки GET для получения заказа по ID
func TestGetOrderByID(t *testing.T) {
	gin.SetMode(gin.TestMode) // чтобы не было лишних логов

	// загружаем данные из файлов перед тестом
	if err := repository.LoadOrdersFromFile("../../data/orders.json"); err != nil {
		t.Fatalf("Не удалось загрузить заказы: %v", err)
	}
	if err := repository.LoadUsersFromFile("../../data/users.json"); err != nil {
		t.Fatalf("Не удалось загрузить пользователей: %v", err)
	}

	tests := []struct {
		name       string
		orderID    string
		wantStatus int
		wantID     string
	}{
		{
			name:       "existing order",
			orderID:    "Order-1753281561910859000",
			wantStatus: http.StatusOK,
			wantID:     "Order-1753281561910859000",
		},
		{
			name:       "non-existing order",
			orderID:    "non-existent-id",
			wantStatus: http.StatusNotFound,
			wantID:     "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// создаем сервер
			s := NewServer(":8080")

			// получаем роутер
			r := s.httpServer.Handler.(*gin.Engine)

			// создаем запрос
			req, _ := http.NewRequest("GET", "/api/orders/"+tc.orderID, nil)
			w := httptest.NewRecorder()

			// выполняем запрос
			r.ServeHTTP(w, req)

			// проверяем код ответа
			assert.Equal(t, tc.wantStatus, w.Code)

			// проверяем тело ответа
			if tc.wantStatus == http.StatusOK {
				var order model.Order
				err := json.Unmarshal(w.Body.Bytes(), &order)
				assert.NoError(t, err)
				assert.Equal(t, tc.orderID, order.Id)
			}
		})
	}
}

func TestCreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s := NewServer(":8080")
	r := s.httpServer.Handler.(*gin.Engine)

	tests := []struct {
		name        string
		body        string
		wantStatus  int
		wantUserID  string
		wantCreated bool
	}{
		{
			name:        "valid order",
			body:        `{"user_id":"User-testOne", "status":0}`,
			wantStatus:  http.StatusCreated,
			wantUserID:  "User-testOne",
			wantCreated: true,
		},
		{
			name:       "missing user_id",
			body:       `{"status":1}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       `{user_id:"User-testOne", "status":}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/orders", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			if tc.wantCreated {
				var got model.Order
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.NotEmpty(t, got.Id)
				assert.Equal(t, tc.wantUserID, got.UserID)
			}
		})
	}
}

func TestDeleteOrderByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s := NewServer(":8080")
	r := s.httpServer.Handler.(*gin.Engine)

	tests := []struct {
		name       string
		prepare    func() string // возвращает ID заказа для удаления
		wantStatus int
	}{
		{
			name: "delete existing order",
			prepare: func() string {
				// создаём новый заказ
				order := model.NewOrder("user-test-delete")
				repository.SaveStorable(order)
				return order.Id
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "delete non-existing order",
			prepare: func() string {
				// возвращаем ID, которого точно нет
				return "non-existent-id"
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			orderID := tc.prepare()

			req, _ := http.NewRequest("DELETE", "/api/orders/"+orderID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

// тест для обновления статуса заказа
func TestOrderStatusHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// создаём тестовые заказы
	orderCreated := model.NewOrder("User1") // статус OrderCreated
	orderConfirmed := model.NewOrder("User2")
	orderConfirmed.Status = model.OrderConfirmed // статус Confirmed
	orderDelivered := model.NewOrder("User3")
	orderDelivered.Status = model.OrderDelivered // статус Delivered
	orderCancelled := model.NewOrder("User4")
	orderCancelled.Status = model.OrderCancelled // статус Cancelled

	// сохраняем в репозиторий
	repository.SaveStorable(orderCreated)
	repository.SaveStorable(orderConfirmed)
	repository.SaveStorable(orderDelivered)
	repository.SaveStorable(orderCancelled)

	s := NewServer(":8080")
	r := s.httpServer.Handler.(*gin.Engine)

	tests := []struct {
		name           string
		route          string
		orderID        string
		wantHTTPStatus int
		wantRepoStatus model.OrderStatus
	}{
		{
			name:           "confirm created order",
			route:          "/api/orders/confirm/",
			orderID:        orderCreated.Id,
			wantHTTPStatus: http.StatusOK,
			wantRepoStatus: model.OrderConfirmed,
		},
		{
			name:           "confirm already confirmed order",
			route:          "/api/orders/confirm/",
			orderID:        orderConfirmed.Id,
			wantHTTPStatus: http.StatusConflict,
		},
		{
			name:           "delivery confirmed order",
			route:          "/api/orders/delivery/",
			orderID:        orderConfirmed.Id,
			wantHTTPStatus: http.StatusOK,
			wantRepoStatus: model.OrderDelivered,
		},
		{
			name:           "cancel created order",
			route:          "/api/orders/cancel/",
			orderID:        orderCreated.Id,
			wantHTTPStatus: http.StatusNoContent,
		},
		{
			name:           "cancel delivered order",
			route:          "/api/orders/cancel/",
			orderID:        orderConfirmed.Id,
			wantHTTPStatus: http.StatusConflict,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", tc.route+tc.orderID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			assert.Equal(t, tc.wantHTTPStatus, w.Code)

			// проверяем статус в репозитории, если указан
			if tc.wantRepoStatus != 0 {
				order := repository.GetOrderByID(tc.orderID)
				assert.NotNil(t, order)
				assert.Equal(t, tc.wantRepoStatus, order.Status)
			}
		})
	}
}

// тест ручки GET для получения пользователей
func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode) // чтобы не было лишних логов

	if err := repository.LoadUsersFromFile("../../data/users.json"); err != nil {
		t.Fatalf("Не удалось загрузить пользователей: %v", err)
	}

	tests := []struct {
		name        string
		expectedMin int
	}{
		{
			name:        "users exist",
			expectedMin: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// создаем сервер
			s := NewServer(":8080")

			// получаем роутер
			r := s.httpServer.Handler.(*gin.Engine)

			// создаем запрос
			req, _ := http.NewRequest("GET", "/api/users", nil)
			w := httptest.NewRecorder()

			// выполняем запрос
			r.ServeHTTP(w, req)

			// проверяем код ответа
			assert.Equal(t, http.StatusOK, w.Code)

			// проверяем тело ответа
			var got []model.User
			err := json.Unmarshal(w.Body.Bytes(), &got)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(got), tc.expectedMin)
		})
	}
}

// тест ручки GET для получения пользователя по ID
func TestGetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode) // чтобы не было лишних логов

	if err := repository.LoadUsersFromFile("../../data/users.json"); err != nil {
		t.Fatalf("Не удалось загрузить пользователей: %v", err)
	}

	tests := []struct {
		name       string
		userID     string
		wantStatus int
		wantID     string
	}{
		{
			name:       "existing user",
			userID:     "User-1754393816395113000",
			wantStatus: http.StatusOK,
			wantID:     "User-1754393816395113000",
		},
		{
			name:       "non-existing user",
			userID:     "non-existent-id",
			wantStatus: http.StatusNotFound,
			wantID:     "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// создаем сервер
			s := NewServer(":8080")

			// получаем роутер
			r := s.httpServer.Handler.(*gin.Engine)

			// создаем запрос
			req, _ := http.NewRequest("GET", "/api/users/"+tc.userID, nil)
			w := httptest.NewRecorder()

			// выполняем запрос
			r.ServeHTTP(w, req)

			// проверяем код ответа
			assert.Equal(t, tc.wantStatus, w.Code)

			// проверяем тело ответа
			if tc.wantStatus == http.StatusOK {
				var user model.User
				err := json.Unmarshal(w.Body.Bytes(), &user)
				assert.NoError(t, err)
				assert.Equal(t, tc.userID, user.Id)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s := NewServer(":8080")
	r := s.httpServer.Handler.(*gin.Engine)

	tests := []struct {
		name        string
		body        string
		wantStatus  int
		wantCreated bool
	}{
		{
			name:        "valid user",
			body:        `{"id":"User-test", "name":"Гера"}`,
			wantStatus:  http.StatusCreated,
			wantCreated: true,
		},
		{
			name:       "invalid json",
			body:       `{id:"User-test", "name": ""}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/users", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			if tc.wantCreated {
				var got model.User
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.NotEmpty(t, got.Id)
			}
		})
	}
}

func TestDeleteUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s := NewServer(":8080")
	r := s.httpServer.Handler.(*gin.Engine)

	tests := []struct {
		name       string
		prepare    func() string // возвращает ID заказа для удаления
		wantStatus int
	}{
		{
			name: "delete existing user",
			prepare: func() string {
				// создаём новый заказ
				user := model.NewUser("user-test1-delete")
				repository.SaveStorable(user)
				return user.Id
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "delete non-existing user",
			prepare: func() string {
				// возвращаем ID, которого точно нет
				return "non-existent-id"
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userID := tc.prepare()

			req, _ := http.NewRequest("DELETE", "/api/users/"+userID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

func TestUserUpdateByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s := NewServer(":8080")
	r := s.httpServer.Handler.(*gin.Engine)

	// создаём пользователя для тестов
	existingUserID := "u1"
	repository.SaveStorable(&model.User{Id: existingUserID, Name: "Old Name"})

	tests := []struct {
		name           string
		userID         string
		body           any
		expectedStatus int
		expectedName   string
	}{
		{
			name:           "успешное обновление существующего пользователя",
			userID:         existingUserID,
			body:           updateUserRequest{Name: "New Name"},
			expectedStatus: http.StatusOK,
			expectedName:   "New Name",
		},
		{
			name:           "пустое имя в запросе",
			userID:         existingUserID,
			body:           updateUserRequest{Name: ""},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "несуществующий пользователь",
			userID:         "not-exist",
			body:           updateUserRequest{Name: "Ghost"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "невалидный JSON",
			userID:         existingUserID,
			body:           "{invalid-json", // специально строка вместо структуры
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			switch v := tc.body.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(v)
				assert.NoError(t, err)
			}

			req, _ := http.NewRequest(http.MethodPut, "/api/users/"+tc.userID, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedStatus == http.StatusOK {
				var updatedUser model.User
				err := json.Unmarshal(w.Body.Bytes(), &updatedUser)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedName, updatedUser.Name)

				user := repository.GetUserByID(tc.userID)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectedName, user.Name)
			}
		})
	}
}

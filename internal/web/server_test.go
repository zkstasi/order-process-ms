package web

import (
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

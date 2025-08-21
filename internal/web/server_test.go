package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"order-ms/internal/model"
	"order-ms/internal/repository"
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

package service_test

import (
	"order-ms/internal/model"
	"order-ms/internal/service"
	"testing"
)

// MockRepo — простая реализация интерфейса service.Repository для тестов
type MockRepo struct {
	Saved      []model.Storable
	Orders     []model.Order
	Users      []model.User
	Deliveries []model.Delivery
	Warehouses []model.Warehouse
}

// методы, возвращают заранее подготовленные данные/ сохраняет их в память

func (m *MockRepo) Save(s model.Storable) error {
	m.Saved = append(m.Saved, s)
	return nil
}

func (m *MockRepo) GetOrders() []model.Order {
	return m.Orders
}

func (m *MockRepo) GetUsers() []model.User {
	return m.Users
}

func (m *MockRepo) GetDeliveries() []model.Delivery {
	return m.Deliveries
}

func (m *MockRepo) GetWarehouses() []model.Warehouse {
	return m.Warehouses
}

// Тест

func TestProcessDataChan(t *testing.T) {

	// срез структур - таблица тестов
	tests := []struct {
		name     string           // имя кейса
		inputs   []model.Storable // набор объектов, которые кладем в канал
		expected int              // кол-во сохраненных объектов, которое ожидаем увидеть
	}{
		// первый сценарий
		{
			name:     "user-1",
			inputs:   []model.Storable{model.NewUser("Иван")},
			expected: 1,
		},

		// второй сценарий
		{
			name:     "user and order",
			inputs:   []model.Storable{model.NewUser("Маша"), model.NewOrder("some-user-id")},
			expected: 2,
		},
	}

	// итерируемся по всем тест-кейсам
	for _, tc := range tests {
		//создаем саб-тест
		t.Run(tc.name, func(t *testing.T) {
			//поднимаем новый мок-репозиторий
			mock := &MockRepo{}
			// создаем буферизированный канал на len(tc.inputs) элементов
			dataChan := make(chan model.Storable, len(tc.inputs))

			// отправляем все входные объекты из таблицы в канал
			for _, s := range tc.inputs {
				dataChan <- s
			}
			close(dataChan)

			service.ProcessDataChan(dataChan, mock) // вызов тестируемой функции

			// проверяем сколько объектов сохранил мок
			if len(mock.Saved) != tc.expected {
				t.Errorf("Expected %d saved objects, got %d", tc.expected, len(mock.Saved))
			}
		})
	}
}

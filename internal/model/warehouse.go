package model

import "time"

type WarehouseStatus int

type Warehouse struct {
	Id      int64           `json:"id"`      // Уникальный идентификатор склада
	OrderId string          `json:"orderId"` // ID заказа
	Status  WarehouseStatus `json:"status"`
}

// NewWarehouse создаёт новый склад с заданным id, заказом и статусом.

func NewWarehouse(orderId string, status WarehouseStatus) *Warehouse {
	return &Warehouse{
		Id:      generateWarehouseID(),
		OrderId: orderId,
		Status:  status,
	}
}

func generateWarehouseID() int64 {
	return time.Now().UnixNano()
}

// реализация интерфейса Storable

func (w *Warehouse) GetType() string {
	return "warehouse"
}

package model

type WarehouseStatus int

type Warehouse struct {
	id      int // Уникальный идентификатор склада
	orderId int // ID заказа
	status  WarehouseStatus
}

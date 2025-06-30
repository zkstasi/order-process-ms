package model

type WarehouseStatus int

type Warehouse struct {
	id      int // Уникальный идентификатор склада
	orderId int // ID заказа
	status  WarehouseStatus
}

func (w *Warehouse) Id() int {
	return w.id
}

func (w *Warehouse) OrderId() int {
	return w.orderId
}

func (w *Warehouse) Status() WarehouseStatus {
	return w.status
}

func (w *Warehouse) SetStatus(newStatus WarehouseStatus) {
	w.status = newStatus
}

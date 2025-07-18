package model

type WarehouseStatus int

type Warehouse struct {
	id      int    // Уникальный идентификатор склада
	orderId string // ID заказа
	status  WarehouseStatus
}

// NewWarehouse создаёт новый склад с заданным id, заказом и статусом.

func NewWarehouse(id int, orderId string, status WarehouseStatus) *Warehouse {
	return &Warehouse{
		id:      id,
		orderId: orderId,
		status:  status,
	}
}

func (w *Warehouse) Id() int {
	return w.id
}

func (w *Warehouse) OrderId() string {
	return w.orderId
}

func (w *Warehouse) Status() WarehouseStatus {
	return w.status
}

func (w *Warehouse) SetStatus(newStatus WarehouseStatus) {
	w.status = newStatus
}

// реализация интерфейса Storable

func (w *Warehouse) GetType() string {
	return "warehouse"
}

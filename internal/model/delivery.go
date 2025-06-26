package model

type DeliveryStatus int

type Delivery struct {
	id      int    // Уникальный идентификатор доставки
	orderId int    // ID заказа
	userId  int    // ID клиента
	address string // Адрес доставки
	status  DeliveryStatus
}

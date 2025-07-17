package model

import (
	"time"
)

type DeliveryStatus int

type Delivery struct {
	Id      int64          `json:"id"`      // Уникальный идентификатор доставки
	OrderId string         `json:"OrderId"` // ID заказа
	UserId  string         `json:"UserId"`  // ID клиента
	Address string         `json:"Address"` // Адрес доставки
	Status  DeliveryStatus `json:"status"`  // Статус доставки
}

// NewDelivery создаёт новую доставку с заданными параметрами.
// Статус доставки устанавливается по умолчанию в 0 ("новая").

func NewDelivery(orderId string, userId string, address string, status DeliveryStatus) *Delivery {
	return &Delivery{
		Id:      generateDeliveryID(),
		OrderId: orderId,
		UserId:  userId,
		Address: address,
		Status:  status,
	}
}

func generateDeliveryID() int64 {
	return time.Now().UnixNano()
}

// реализация интерфейса Storable

func (d *Delivery) GetType() string {
	return "delivery"
}

package model

import (
	"fmt"
	"time"
)

// создание нового типа для статуса заказа

type OrderStatus int

// структура для объекта Заказ

type Order struct {
	Id        string      `json:"id"`         // Уникальный номер заказа
	UserID    string      `json:"user_id"`    // Кто сделал заказ
	Status    OrderStatus `json:"status"`     // Статус заказа (0-3)
	CreatedAt time.Time   `json:"created_at"` // Когда заказ создан
}

// NewOrder создаёт новый заказ с уникальным ID, привязанный к пользователю userID.
// Статус по умолчанию — 0 (новый заказ).

func NewOrder(newUserId string) *Order {
	return &Order{
		Id:        generateUniqID(),
		UserID:    newUserId,
		Status:    OrderStatus(0),
		CreatedAt: time.Now(),
	}
}

// функция для генерации id заказа

func generateUniqID() string {
	return fmt.Sprintf("Order-%d", time.Now().UnixNano())
}

// реализация интерфейса Storable

func (o *Order) GetType() string {
	return "order"
}

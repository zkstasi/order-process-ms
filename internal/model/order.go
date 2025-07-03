package model

import (
	"fmt"
	"time"
)

// создание нового типа для статуса заказа

type OrderStatus int

//const (
//	OrderCreated OrderStatus = iota
//	OrderAborted
//)

// структура для объекта Заказ

type Order struct {
	id        string      // Приватный уникальный номер заказа
	userID    string      // Кто сделал заказ
	status    OrderStatus // Приватный статус заказа (0-3)
	createdAt time.Time   // Когда заказ создан
}

// NewOrder создаёт новый заказ с уникальным ID, привязанный к пользователю userID.
// Статус по умолчанию — 0 (новый заказ).

func NewOrder(newUserId string) *Order {
	return &Order{
		id:        generateUniqID(),
		userID:    newUserId,
		status:    OrderStatus(0),
		createdAt: time.Now(),
	}
}

// методы для получения полей структуры

func (o *Order) Id() string {
	return o.id
}

func (o *Order) UserId() string {
	return o.userID
}

func (o *Order) Status() OrderStatus {
	return o.status
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// метод для изменения поля статус

func (o *Order) SetStatus(newStatus OrderStatus) {
	o.status = newStatus
}

// функция для генерации id заказа

func generateUniqID() string {
	return fmt.Sprintf("Order-%d", time.Now().UnixNano())
}

// реализация интерфейса Storable

func (o *Order) GetType() string {
	return "order"
}

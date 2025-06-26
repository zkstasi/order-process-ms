package model

import (
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

// функция для создания готового заказа, вместо SetAll

func NewOrder(newUserID string) *Order {
	return &Order{
		id:        generateUniqID(),
		userID:    newUserID,
		status:    OrderStatus(0),
		createdAt: time.Now(),
	}
}

// методы для получения полей структуры

func (o *Order) ID() string {
	return o.id
}

func (o *Order) UserID() string {
	return o.userID
}

func (o *Order) Status() OrderStatus {
	return o.status
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) Fields() (string, string, OrderStatus, time.Time) {
	return o.id, o.userID, o.status, o.createdAt
}

// метод для изменения поля статус

func (o *Order) SetStatus(newStatus OrderStatus) {
	o.status = newStatus
}

// функция для генерации id заказа

func generateUniqID() string {
	return "" // тут должна быть логика генерации id заказа
}

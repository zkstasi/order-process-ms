package order

import "time"

type Order struct {
	id        int64     // Приватный уникальный номер
	userID    string    // Кто сделал заказ
	status    int       // Приватный статус заказа (0-3)
	createdAt time.Time // Когда заказ создан
}

func (o *Order) GetID() int64 {
	return o.id
}

func (o *Order) GetUserID() string {
	return o.userID
}

func (o *Order) GetStatus() int {
	return o.status
}

func (o *Order) GetCreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) GetAll() (int64, string, int, time.Time) {
	return o.id, o.userID, o.status, o.createdAt
}

func (o *Order) SetID(newID int64) {
	o.id = newID
}

func (o *Order) SetUserID(newUserID string) {
	o.userID = newUserID
}

func (o *Order) SetStatus(newStatus int) {
	o.status = newStatus
}

func (o *Order) SetAll(newUserID string, newStatus int) {
	o.userID = newUserID
	o.status = newStatus
}

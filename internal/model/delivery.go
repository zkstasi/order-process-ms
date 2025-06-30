package model

type DeliveryStatus int

type Delivery struct {
	id      int    // Уникальный идентификатор доставки
	orderId int    // ID заказа
	userId  int    // ID клиента
	address string // Адрес доставки
	status  DeliveryStatus
}

func (d *Delivery) Id() int {
	return d.id
}

func (d *Delivery) OrderId() int {
	return d.orderId
}

func (d *Delivery) UserId() int {
	return d.userId
}

func (d *Delivery) Address() string {
	return d.address
}

func (d *Delivery) Status() DeliveryStatus {
	return d.status
}

func (d *Delivery) SetStatus(newStatus DeliveryStatus) {
	d.status = newStatus

}

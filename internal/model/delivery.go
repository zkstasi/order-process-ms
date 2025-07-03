package model

type DeliveryStatus int

type Delivery struct {
	id      int    // Уникальный идентификатор доставки
	orderId string // ID заказа
	userId  string // ID клиента
	address string // Адрес доставки
	status  DeliveryStatus
}

// NewDelivery создаёт новую доставку с заданными параметрами.
// Статус доставки устанавливается по умолчанию в 0 ("новая").

func NewDelivery(id int, orderId string, userId string, address string, status DeliveryStatus) *Delivery {
	return &Delivery{
		id:      id,
		orderId: orderId,
		userId:  userId,
		address: address,
		status:  status,
	}
}

func (d *Delivery) Id() int {
	return d.id
}

func (d *Delivery) OrderId() string {
	return d.orderId
}

func (d *Delivery) UserId() string {
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

// реализация интерфейса Storable

func (d *Delivery) GetType() string {
	return "delivery"
}

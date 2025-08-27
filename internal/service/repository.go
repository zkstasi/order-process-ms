package service

import "order-ms/internal/model"

type Repository interface {
	// Общий Save
	Save(s model.Storable) error

	// Заказы
	SaveOrder(order *model.Order) error
	GetOrders() ([]*model.Order, error)
	GetOrderByID(id string) (*model.Order, error)
	DeleteOrder(id string) (bool, error)
	ConfirmOrder(orderId string) (bool, error)
	DeliverOrder(id string) (bool, error)
	CancelOrder(id string) (bool, error)

	// Пользователи
	SaveUser(user *model.User) error
	GetUsers() ([]*model.User, error)
	GetUserByID(id string) (*model.User, error)
	UpdateUserName(id, name string) (bool, error)
	DeleteUser(id string) (bool, error)

	// доставки и склады
	GetDeliveries() ([]*model.Delivery, error)
	GetWarehouses() ([]*model.Warehouse, error)
}

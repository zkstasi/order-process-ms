package service

import "order-ms/internal/model"

type Repository interface {
	Save(s model.Storable) error
	GetOrders() []model.Order
	GetUsers() []model.User
	GetDeliveries() []model.Delivery
	GetWarehouses() []model.Warehouse
}

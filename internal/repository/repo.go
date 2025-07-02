package repository

import (
	"fmt"
	"order-ms/internal/model"
)

var orders []*model.Order //слайс структуры Order
var users []*model.User
var deliveries []*model.Delivery
var warehouses []*model.Warehouse

//функция, принимает любой объект, реализующий интерфейс
//проверяет конкретный тип и добавляет его в соответствующий слайс

func SaveStorable(s model.Storable) {
	switch v := s.(type) {
	case *model.Order:
		orders = append(orders, v)
	case *model.User:
		users = append(users, v)
	case *model.Delivery:
		deliveries = append(deliveries, v)
	case *model.Warehouse:
		warehouses = append(warehouses, v)
	default:
		fmt.Println("неизвестный тип")
	}
}

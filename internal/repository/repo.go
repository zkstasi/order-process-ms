package repository

import (
	"fmt"
	"order-ms/internal/model"
	"sync"
)

var (
	orders     []*model.Order //слайс структуры Order
	users      []*model.User
	deliveries []*model.Delivery
	warehouses []*model.Warehouse

	muOrders     sync.Mutex // Защита слайс от гонок данных
	muUsers      sync.Mutex
	muDeliveries sync.Mutex
	muWarehouses sync.Mutex
)

//функция, принимает любой объект, реализующий интерфейс
//проверяет конкретный тип и добавляет его в соответствующий слайс

func SaveStorable(dataChan <-chan model.Storable) {
	for s := range dataChan {
		switch v := s.(type) {
		case *model.Order:
			muOrders.Lock()
			orders = append(orders, v)
			muOrders.Unlock()
		case *model.User:
			muUsers.Lock()
			users = append(users, v)
			muUsers.Unlock()
		case *model.Delivery:
			muDeliveries.Lock()
			deliveries = append(deliveries, v)
			muDeliveries.Unlock()
		case *model.Warehouse:
			muWarehouses.Lock()
			warehouses = append(warehouses, v)
			muWarehouses.Unlock()
		default:
			fmt.Println("Type: Undefined")
		}
	}

}

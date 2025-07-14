package repository

import (
	"fmt"
	"order-ms/internal/model"
	"sync"
)

var (
	Orders     []*model.Order //слайс структуры Order
	Users      []*model.User
	Deliveries []*model.Delivery
	Warehouses []*model.Warehouse

	MuOrders     sync.Mutex // Защита слайс от гонок данных
	MuUsers      sync.Mutex
	MuDeliveries sync.Mutex
	MuWarehouses sync.Mutex
)

//функция, принимает любой объект, реализующий интерфейс
//проверяет конкретный тип и добавляет его в соответствующий слайс

func SaveStorable(dataChan <-chan model.Storable) {
	for s := range dataChan {
		switch v := s.(type) {
		case *model.Order:
			MuOrders.Lock()
			Orders = append(Orders, v)
			MuOrders.Unlock()
		case *model.User:
			MuUsers.Lock()
			Users = append(Users, v)
			MuUsers.Unlock()
		case *model.Delivery:
			MuDeliveries.Lock()
			Deliveries = append(Deliveries, v)
			MuDeliveries.Unlock()
		case *model.Warehouse:
			MuWarehouses.Lock()
			Warehouses = append(Warehouses, v)
			MuWarehouses.Unlock()
		default:
			fmt.Println("Type: Undefined")
		}
	}

}

package repository

import (
	"fmt"
	"order-ms/internal/model"
	"sync"
	"time"
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

func Logger(stop <-chan struct{}) {

	// предыдущие длины слайс

	lastOrdersCount := 0
	lastUsersCount := 0
	lastDeliveriesCount := 0
	lastWarehousesCount := 0

	for {
		select {
		case <-stop:
			return
		default:
			muOrders.Lock()
			ordersCount := len(orders)
			muOrders.Unlock()

			muUsers.Lock()
			usersCount := len(users)
			muUsers.Unlock()

			muDeliveries.Lock()
			deliveriesCount := len(deliveries)
			muDeliveries.Unlock()

			muWarehouses.Lock()
			warehousesCount := len(warehouses)
			muWarehouses.Unlock()

			if ordersCount > lastOrdersCount {
				fmt.Printf("New orders: %d\n", ordersCount-lastOrdersCount)
				lastOrdersCount = ordersCount
			}
			if usersCount > lastUsersCount {
				fmt.Printf("New users: %d\n", usersCount-lastUsersCount)
				lastUsersCount = usersCount
			}
			if deliveriesCount > lastDeliveriesCount {
				fmt.Printf("New deliveries: %d\n", deliveriesCount-lastDeliveriesCount)
				lastDeliveriesCount = deliveriesCount
			}
			if warehousesCount > lastWarehousesCount {
				fmt.Printf("New warehouses: %d\n", warehousesCount-lastWarehousesCount)
				lastWarehousesCount = warehousesCount
			}
		}

		time.Sleep(200 * time.Millisecond)
	}
}

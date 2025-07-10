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

	lastOrdersIndex := 0
	lastUsersIndex := 0
	lastDeliveriesIndex := 0
	lastWarehousesIndex := 0

	for {
		select {
		case <-stop:
			return
		default:
			muOrders.Lock()
			ordersCount := len(orders)
			newOrders := orders[lastOrdersIndex:ordersCount]
			muOrders.Unlock()

			muUsers.Lock()
			usersCount := len(users)
			newUsers := users[lastUsersIndex:usersCount]
			muUsers.Unlock()

			muDeliveries.Lock()
			deliveriesCount := len(deliveries)
			newDeliveries := deliveries[lastDeliveriesIndex:deliveriesCount]
			muDeliveries.Unlock()

			muWarehouses.Lock()
			warehousesCount := len(warehouses)
			newWarehouses := warehouses[lastWarehousesIndex:warehousesCount]
			muWarehouses.Unlock()

			if len(newOrders) > 0 {
				fmt.Printf("New orders: %d\n", len(newOrders))
				for _, o := range newOrders {
					fmt.Printf("Order ID: %s, UserID: %s, Status: %d, CreatedAt: %s\n", o.Id(), o.UserId(), o.Status(), o.CreatedAt())
				}
				lastOrdersIndex = ordersCount
			}
			if len(newUsers) > 0 {
				fmt.Printf("New users: %d\n", len(newUsers))
				for _, u := range newUsers {
					fmt.Printf("User ID: %s, Name: %s\n", u.Id(), u.Name())
				}
				lastUsersIndex = usersCount
			}

			if len(newDeliveries) > 0 {
				fmt.Printf("New deliveries: %d\n", len(newDeliveries))
				for _, d := range newDeliveries {
					fmt.Printf("Delivery ID: %d, OrderID: %s, UserID: %s, Address: %s, Status: %d\n", d.Id(), d.OrderId(), d.UserId(), d.Address(), d.Status())
				}
				lastDeliveriesIndex = deliveriesCount
			}

			if len(newWarehouses) > 0 {
				fmt.Printf("New warehouses: %d\n", len(newWarehouses))
				for _, w := range newWarehouses {
					fmt.Printf("Warehouse ID: %d, OrderID: %s, Status: %d\n", w.Id(), w.OrderId(), w.Status())
				}
				lastWarehousesIndex = warehousesCount
			}
		}

		time.Sleep(200 * time.Millisecond)
	}
}

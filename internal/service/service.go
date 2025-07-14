package service

import (
	"fmt"
	"order-ms/internal/model"
	"order-ms/internal/repository"
	"time"
)

// функция для создания структур и передачи их в репозиторий через канал

func CreateStructs(dataChan chan<- model.Storable, stop <-chan struct{}) {

	for {
		select {
		case <-stop: // Остановка бесконечного цикла
			return
		default:
			order := model.NewOrder("user123")
			dataChan <- order

			user := model.NewUser("user123", "Петя")
			dataChan <- user

			delivery := model.NewDelivery(65, "order-783", "user123", "ул. Ленина", 0)
			dataChan <- delivery

			warehouse := model.NewWarehouse(543, "order-783", 0)
			dataChan <- warehouse

			time.Sleep(300 * time.Millisecond) // Пауза между отправками
		}

	}
}

func Logger(loggerStop <-chan struct{}) {

	// предыдущие длины слайс

	lastOrdersIndex := 0
	lastUsersIndex := 0
	lastDeliveriesIndex := 0
	lastWarehousesIndex := 0

	for {
		select {
		case <-loggerStop:
			return
		default:
			repository.MuOrders.Lock()
			ordersCount := len(repository.Orders)
			newOrders := repository.Orders[lastOrdersIndex:ordersCount]
			repository.MuOrders.Unlock()

			repository.MuUsers.Lock()
			usersCount := len(repository.Users)
			newUsers := repository.Users[lastUsersIndex:usersCount]
			repository.MuUsers.Unlock()

			repository.MuDeliveries.Lock()
			deliveriesCount := len(repository.Deliveries)
			newDeliveries := repository.Deliveries[lastDeliveriesIndex:deliveriesCount]
			repository.MuDeliveries.Unlock()

			repository.MuWarehouses.Lock()
			warehousesCount := len(repository.Warehouses)
			newWarehouses := repository.Warehouses[lastWarehousesIndex:warehousesCount]
			repository.MuWarehouses.Unlock()

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

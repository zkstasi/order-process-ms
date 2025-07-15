package service

import (
	"fmt"
	"order-ms/internal/model"
	"order-ms/internal/repository"
	"time"
)

// функция для создания структур и передачи их в канал DataChan

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

// функция, которая читает из DataChan и сохраняет данные в репозиторий

func ProcessDataChan(dataChan <-chan model.Storable, stop <-chan struct{}) {
	for {
		select {
		case s, ok := <-dataChan:
			if !ok {
				return
			}
			repository.SaveStorable(s)
		case <-stop:
			return
		}
	}
}

func Logger(loggerStop <-chan struct{}) {

	// получаем стартовые длины, чтобы считать только новые данные

	lastOrdersIndex := len(repository.GetOrders())
	lastUsersIndex := len(repository.GetUsers())
	lastDeliveriesIndex := len(repository.GetDeliveries())
	lastWarehousesIndex := len(repository.GetWarehouses())

	for {
		select {
		case <-loggerStop:
			return
		default:
			orders := repository.GetOrders() // вызов функции, возвращаем копию среза и сохраняем в переменную
			ordersCount := len(orders)
			newOrders := orders[lastOrdersIndex:ordersCount]

			if len(newOrders) > 0 {
				fmt.Printf("New orders: %d\n", len(newOrders))
				for _, o := range newOrders {
					fmt.Printf("Order ID: %s, UserID: %s, Status: %d, CreatedAt: %s\n", o.Id(), o.UserId(), o.Status(), o.CreatedAt())
				}
				lastOrdersIndex = ordersCount
			}

			users := repository.GetUsers()
			usersCount := len(users)
			newUsers := users[lastUsersIndex:usersCount]

			if len(newUsers) > 0 {
				fmt.Printf("New users: %d\n", len(newUsers))
				for _, u := range newUsers {
					fmt.Printf("User ID: %s, Name: %s\n", u.Id(), u.Name())
				}
				lastUsersIndex = usersCount
			}

			deliveries := repository.GetDeliveries()
			deliveriesCount := len(deliveries)
			newDeliveries := deliveries[lastDeliveriesIndex:deliveriesCount]

			if len(newDeliveries) > 0 {
				fmt.Printf("New deliveries: %d\n", len(newDeliveries))
				for _, d := range newDeliveries {
					fmt.Printf("Delivery ID: %d, OrderID: %s, UserID: %s, Address: %s, Status: %d\n", d.Id(), d.OrderId(), d.UserId(), d.Address(), d.Status())
				}
				lastDeliveriesIndex = deliveriesCount
			}

			warehouses := repository.GetWarehouses()
			warehousesCount := len(warehouses)
			newWarehouses := warehouses[lastWarehousesIndex:warehousesCount]

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

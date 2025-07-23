package service

import (
	"context"
	"fmt"
	"order-ms/internal/model"
	"order-ms/internal/repository"
	"time"
)

// функция для создания структур и передачи их в канал DataChan

func CreateStructs(ctx context.Context, dataChan chan<- model.Storable) {
	for {
		select {
		case <-ctx.Done(): // Контекст отменен, нужно завершить работу
			return
		case <-time.After(200 * time.Millisecond): // контекст отменяется без искусственной задержки

			// создаем пользователя
			user := model.NewUser("Петя")
			dataChan <- user

			// создаем заказ с id пользователя
			order := model.NewOrder(user.Id)
			dataChan <- order

			// создаем доставку с id заказа и пользователя
			delivery := model.NewDelivery(order.Id, user.Id, "ул. Ленина", 0)
			dataChan <- delivery

			// создаем склад с id заказа
			warehouse := model.NewWarehouse(order.Id, 0)
			dataChan <- warehouse
		}
	}
}

// функция, которая читает из DataChan и сохраняет данные в репозиторий через бесконечный цикл
// его останавливает закрытие канала dataChan

func ProcessDataChan(dataChan <-chan model.Storable) {
	for s := range dataChan {
		repository.SaveStorable(s)
	}
}

func Logger(ctx context.Context) {

	// получаем стартовые длины, чтобы считать только новые данные

	lastOrdersIndex := len(repository.GetOrders())
	lastUsersIndex := len(repository.GetUsers())
	lastDeliveriesIndex := len(repository.GetDeliveries())
	lastWarehousesIndex := len(repository.GetWarehouses())

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(200 * time.Millisecond):
			orders := repository.GetOrders() // вызов функции, возвращаем копию среза и сохраняем в переменную
			ordersCount := len(orders)
			if lastOrdersIndex > ordersCount {
				lastOrdersIndex = ordersCount
			}
			newOrders := orders[lastOrdersIndex:ordersCount]

			if len(newOrders) > 0 {
				fmt.Printf("New orders: %d\n", len(newOrders))
				for _, o := range newOrders {
					fmt.Printf("Order ID: %s, UserID: %s, Status: %d, CreatedAt: %s\n", o.Id, o.UserID, o.Status, o.CreatedAt)
				}
				lastOrdersIndex = ordersCount
			}

			users := repository.GetUsers()
			usersCount := len(users)
			if lastUsersIndex > usersCount {
				lastUsersIndex = usersCount
			}
			newUsers := users[lastUsersIndex:usersCount]

			if len(newUsers) > 0 {
				fmt.Printf("New users: %d\n", len(newUsers))
				for _, u := range newUsers {
					fmt.Printf("User ID: %s, Name: %s\n", u.Id, u.Name)
				}
				lastUsersIndex = usersCount
			}

			deliveries := repository.GetDeliveries()
			deliveriesCount := len(deliveries)
			if lastDeliveriesIndex > deliveriesCount {
				lastDeliveriesIndex = deliveriesCount
			}
			newDeliveries := deliveries[lastDeliveriesIndex:deliveriesCount]

			if len(newDeliveries) > 0 {
				fmt.Printf("New deliveries: %d\n", len(newDeliveries))
				for _, d := range newDeliveries {
					fmt.Printf("Delivery ID: %d, OrderID: %s, UserID: %s, Address: %s, Status: %d\n", d.Id, d.OrderId, d.UserId, d.Address, d.Status)
				}
				lastDeliveriesIndex = deliveriesCount
			}

			warehouses := repository.GetWarehouses()
			warehousesCount := len(warehouses)
			if lastWarehousesIndex > warehousesCount {
				lastWarehousesIndex = warehousesCount
			}
			newWarehouses := warehouses[lastWarehousesIndex:warehousesCount]

			if len(newWarehouses) > 0 {
				fmt.Printf("New warehouses: %d\n", len(newWarehouses))
				for _, w := range newWarehouses {
					fmt.Printf("Warehouse ID: %d, OrderID: %s, Status: %d\n", w.Id, w.OrderId, w.Status)
				}
				lastWarehousesIndex = warehousesCount
			}
		}
	}
}

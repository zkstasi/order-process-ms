package repository

import (
	"encoding/json"
	"fmt"
	"order-ms/internal/model"
	"os"
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

func SaveStorable(s model.Storable) {
	switch v := s.(type) {
	case *model.Order:
		muOrders.Lock()
		orders = append(orders, v)
		muOrders.Unlock()
		if err := SaveOrdersToFile("data/orders.json"); err != nil {
			fmt.Println("Ошибка при сохранении orders", err)
		}
	case *model.User:
		muUsers.Lock()
		users = append(users, v)
		muUsers.Unlock()
		if err := SaveUsersToFile("data/users.json"); err != nil {
			fmt.Println("Ошибка при сохранении users", err)
		}
	case *model.Delivery:
		muDeliveries.Lock()
		deliveries = append(deliveries, v)
		muDeliveries.Unlock()
		if err := SaveDeliveriesToFile("data/deliveries.json"); err != nil {
			fmt.Println("Ошибка при сохранении deliveries", err)
		}
	case *model.Warehouse:
		muWarehouses.Lock()
		warehouses = append(warehouses, v)
		muWarehouses.Unlock()
		if err := SaveWarehousesToFile("data/warehouses.json"); err != nil {
			fmt.Println("Ошибка при сохранении warehouses", err)
		}
	default:
		fmt.Println("Type: Undefined")
	}
}

// методы получения копий слайсов

func GetOrders() []*model.Order {
	muOrders.Lock()
	defer muOrders.Unlock()

	copiedOrders := make([]*model.Order, len(orders))
	copy(copiedOrders, orders)
	return copiedOrders
}

func GetUsers() []*model.User {
	muUsers.Lock()
	defer muUsers.Unlock()

	copiedUsers := make([]*model.User, len(users))
	copy(copiedUsers, users)
	return copiedUsers
}

func GetDeliveries() []*model.Delivery {
	muDeliveries.Lock()
	defer muDeliveries.Unlock()

	copiedDeliveries := make([]*model.Delivery, len(deliveries))
	copy(copiedDeliveries, deliveries)
	return copiedDeliveries
}

func GetWarehouses() []*model.Warehouse {
	muWarehouses.Lock()
	defer muWarehouses.Unlock()

	copiedWarehouses := make([]*model.Warehouse, len(warehouses))
	copy(copiedWarehouses, warehouses)
	return copiedWarehouses
}

// функции сохранения слайса в json-файл

func SaveOrdersToFile(filepath string) error {
	orders := GetOrders()                             // получаем слайс заказов
	data, err := json.MarshalIndent(orders, "", "  ") // сериализуем в json с отступами и префиксами
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, data, 0644) // сохраняем в указанный файл
	if err != nil {
		return err
	}
	return nil
}

func SaveUsersToFile(filepath string) error {
	users := GetUsers()
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func SaveDeliveriesToFile(filepath string) error {
	deliveries := GetDeliveries()
	data, err := json.MarshalIndent(deliveries, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func SaveWarehousesToFile(filepath string) error {
	warehouses := GetWarehouses()
	data, err := json.MarshalIndent(warehouses, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// функции загрузки json-файлов в слайсы при старте программы

func LoadOrdersFromFile(filepath string) error {
	data, err := os.ReadFile(filepath) // чтение всего файла
	if err != nil {
		return err
	}

	var loadedOrders []model.Order            // временный слайс заказов, так как json не умеет работать с указателями
	err = json.Unmarshal(data, &loadedOrders) // распарсиваем json из data в loadedOrders
	if err != nil {
		return err
	}

	var loadedPointers []*model.Order // слайс, в который поместим указатели на элементы loadedOrders
	for i := range loadedOrders {     // перебираем индексы слайса
		loadedPointers = append(loadedPointers, &loadedOrders[i]) // берем адрес элемента и добавляем указатель в loadedPointers
	}

	muOrders.Lock()
	orders = loadedPointers // заменяем слайс orders на загруженный из файла
	muOrders.Unlock()

	return nil
}

func LoadUsersFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var loadedUsers []model.User
	err = json.Unmarshal(data, &loadedUsers)
	if err != nil {
		return err
	}

	var loadedPointers []*model.User
	for i := range loadedUsers {
		loadedPointers = append(loadedPointers, &loadedUsers[i])
	}

	muUsers.Lock()
	users = loadedPointers
	muUsers.Unlock()

	return nil
}

func LoadDeliveriesFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var loadedDeliveries []model.Delivery
	err = json.Unmarshal(data, &loadedDeliveries)
	if err != nil {
		return err
	}

	var loadedPointers []*model.Delivery
	for i := range loadedDeliveries {
		loadedPointers = append(loadedPointers, &loadedDeliveries[i])
	}

	muDeliveries.Lock()
	deliveries = loadedPointers
	muDeliveries.Unlock()

	return nil
}

func LoadWarehousesFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var loadedWarehouses []model.Warehouse
	err = json.Unmarshal(data, &loadedWarehouses)
	if err != nil {
		return err
	}

	var loadedPointers []*model.Warehouse
	for i := range loadedWarehouses {
		loadedPointers = append(loadedPointers, &loadedWarehouses[i])
	}

	muWarehouses.Lock()
	warehouses = loadedPointers
	muWarehouses.Unlock()

	return nil
}

// функция сохранения данных в файлы

func SaveAllData() {
	err := SaveOrdersToFile("data/orders.json")
	if err != nil {
		fmt.Println("Не удалось сохранить заказы:", err)
	}
	err = SaveUsersToFile("data/users.json")
	if err != nil {
		fmt.Println("Не удалось сохранить пользователей:", err)
	}
	err = SaveDeliveriesToFile("data/deliveries.json")
	if err != nil {
		fmt.Println("Не удалось сохранить доставки:", err)
	}
	err = SaveWarehousesToFile("data/warehouses.json")
	if err != nil {
		fmt.Println("Не удалось сохранить склады:", err)
	}
}

// функция загрузки данных из файлов

func LoadAllData() {
	err := LoadOrdersFromFile("data/orders.json")
	if err != nil {
		fmt.Println("Не удалось загрузить заказы:", err)
	}
	err = LoadUsersFromFile("data/users.json")
	if err != nil {
		fmt.Println("Не удалось загрузить пользователей:", err)
	}
	err = LoadDeliveriesFromFile("data/deliveries.json")
	if err != nil {
		fmt.Println("Не удалось загрузить доставки:", err)
	}
	err = LoadWarehousesFromFile("data/warehouses.json")
	if err != nil {
		fmt.Println("Не удалось загрузить склады:", err)
	}
	fmt.Println("Данные успешно загружены")
}

// метод, который ищет заказ по id

func GetOrderByID(id string) *model.Order {
	muOrders.Lock()
	defer muOrders.Unlock()
	for _, order := range orders {
		if order.Id == id {
			return order
		}
	}
	return nil
}

// методы обновления статуса заказа

func ConfirmOrder(orderId string) bool {
	order := GetOrderByID(orderId) // находим заказ
	if order == nil || order.Status != model.OrderCreated {
		return false
	}
	order.Status = model.OrderConfirmed
	return true
}

func DeliveredOrder(orderId string) bool {
	order := GetOrderByID(orderId)
	if order == nil || order.Status != model.OrderConfirmed {
		return false
	}
	order.Status = model.OrderDelivered
	return true
}

func CancelOrder(orderId string) bool {
	order := GetOrderByID(orderId)
	if order == nil || (order.Status != model.OrderCreated && order.Status != model.OrderConfirmed) {
		return false
	}
	order.Status = model.OrderCancelled
	return true
}

// метод удаления заказа

func DeleteOrder(orderId string) bool {
	muOrders.Lock()
	for i, order := range orders {
		if order.Id == orderId {
			orders = append(orders[:i], orders[i+1:]...)
			muOrders.Unlock()
			if err := SaveOrdersToFile("data/orders.json"); err != nil {
				fmt.Println("Ошибка при сохранении заказов:", err)
			}
			return true
		}
	}
	return false
}

func GetUserByID(id string) *model.User {
	muUsers.Lock()
	defer muUsers.Unlock()
	for _, user := range users {
		if user.Id == id {
			return user
		}
	}
	return nil
}

func UpdateUserName(id, Name string) bool {
	user := GetUserByID(id)
	if user == nil {
		return false
	} else {
		user.Name = Name
	}
	return true
}

func DeleteUser(id string) bool {
	muUsers.Lock()
	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
			muUsers.Unlock()
			if err := SaveUsersToFile("data/users.json"); err != nil {
				fmt.Println("Ошибка при сохранении пользователей:", err)
			}
			return true
		}
	}
	return false
}

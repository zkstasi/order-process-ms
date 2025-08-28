package memory

import (
	"encoding/json"
	"fmt"
	"order-ms/internal/model"
	"os"
	"sync"
)

// Хранит данные в оперативке
type MemoryRepo struct {
	orders     []*model.Order //слайс структуры Order
	users      []*model.User
	deliveries []*model.Delivery
	warehouses []*model.Warehouse

	muOrders     sync.Mutex // Защита слайс от гонок данных
	muUsers      sync.Mutex
	muDeliveries sync.Mutex
	muWarehouses sync.Mutex
}

// конструктор
func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{}
}

//функция, принимает любой объект, реализующий интерфейс
//проверяет конкретный тип и добавляет его в соответствующий слайс

func (r *MemoryRepo) Save(s model.Storable) error {
	switch v := s.(type) {
	case *model.Order:
		r.muOrders.Lock()
		r.orders = append(r.orders, v)
		r.muOrders.Unlock()
		if err := r.SaveOrdersToFile("data/orders.json"); err != nil {
			fmt.Println("Ошибка при сохранении orders", err)
		}
	case *model.User:
		r.muUsers.Lock()
		r.users = append(r.users, v)
		r.muUsers.Unlock()
		if err := r.SaveUsersToFile("data/users.json"); err != nil {
			fmt.Println("Ошибка при сохранении users", err)
		}
	case *model.Delivery:
		r.muDeliveries.Lock()
		r.deliveries = append(r.deliveries, v)
		r.muDeliveries.Unlock()
		if err := r.SaveDeliveriesToFile("data/deliveries.json"); err != nil {
			fmt.Println("Ошибка при сохранении deliveries", err)
		}
	case *model.Warehouse:
		r.muWarehouses.Lock()
		r.warehouses = append(r.warehouses, v)
		r.muWarehouses.Unlock()
		if err := r.SaveWarehousesToFile("data/warehouses.json"); err != nil {
			fmt.Println("Ошибка при сохранении warehouses", err)
		}
	default:
		fmt.Println("Type: Undefined")
	}
	return nil
}

// Обёртки для интерфейса service.Repository
func (r *MemoryRepo) SaveOrder(order *model.Order) error {
	r.Save(order)
	return nil
}

func (r *MemoryRepo) SaveUser(user *model.User) error {
	r.Save(user)
	return nil
}

// методы получения копий слайсов

func (r *MemoryRepo) GetOrders() ([]*model.Order, error) {
	r.muOrders.Lock()
	defer r.muOrders.Unlock()

	copiedOrders := make([]*model.Order, len(r.orders))
	copy(copiedOrders, r.orders)
	return copiedOrders, nil
}

func (r *MemoryRepo) GetUsers() ([]*model.User, error) {
	r.muUsers.Lock()
	defer r.muUsers.Unlock()

	copiedUsers := make([]*model.User, len(r.users))
	copy(copiedUsers, r.users)
	return copiedUsers, nil
}

func (r *MemoryRepo) GetDeliveries() ([]*model.Delivery, error) {
	r.muDeliveries.Lock()
	defer r.muDeliveries.Unlock()

	copiedDeliveries := make([]*model.Delivery, len(r.deliveries))
	copy(copiedDeliveries, r.deliveries)
	return copiedDeliveries, nil
}

func (r *MemoryRepo) GetWarehouses() ([]*model.Warehouse, error) {
	r.muWarehouses.Lock()
	defer r.muWarehouses.Unlock()

	copiedWarehouses := make([]*model.Warehouse, len(r.warehouses))
	copy(copiedWarehouses, r.warehouses)
	return copiedWarehouses, nil
}

// функции сохранения слайса в json-файл

func (r *MemoryRepo) SaveOrdersToFile(filepath string) error {
	orders, _ := r.GetOrders()                        // получаем слайс заказов
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

func (r *MemoryRepo) SaveUsersToFile(filepath string) error {
	users, _ := r.GetUsers()
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

func (r *MemoryRepo) SaveDeliveriesToFile(filepath string) error {
	deliveries, _ := r.GetDeliveries()
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

func (r *MemoryRepo) SaveWarehousesToFile(filepath string) error {
	warehouses, _ := r.GetWarehouses()
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

func (r *MemoryRepo) LoadOrdersFromFile(filepath string) error {
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

	r.muOrders.Lock()
	r.orders = loadedPointers // заменяем слайс orders на загруженный из файла
	r.muOrders.Unlock()

	return nil
}

func (r *MemoryRepo) LoadUsersFromFile(filepath string) error {
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

	r.muUsers.Lock()
	r.users = loadedPointers
	r.muUsers.Unlock()

	return nil
}

func (r *MemoryRepo) LoadDeliveriesFromFile(filepath string) error {
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

	r.muDeliveries.Lock()
	r.deliveries = loadedPointers
	r.muDeliveries.Unlock()

	return nil
}

func (r *MemoryRepo) LoadWarehousesFromFile(filepath string) error {
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

	r.muWarehouses.Lock()
	r.warehouses = loadedPointers
	r.muWarehouses.Unlock()

	return nil
}

// функция сохранения данных в файлы

func (r *MemoryRepo) SaveAllData() {
	err := r.SaveOrdersToFile("data/orders.json")
	if err != nil {
		fmt.Println("Не удалось сохранить заказы:", err)
	}
	err = r.SaveUsersToFile("data/users.json")
	if err != nil {
		fmt.Println("Не удалось сохранить пользователей:", err)
	}
	err = r.SaveDeliveriesToFile("data/deliveries.json")
	if err != nil {
		fmt.Println("Не удалось сохранить доставки:", err)
	}
	err = r.SaveWarehousesToFile("data/warehouses.json")
	if err != nil {
		fmt.Println("Не удалось сохранить склады:", err)
	}
}

// функция загрузки данных из файлов

func (r *MemoryRepo) LoadAllData() {
	err := r.LoadOrdersFromFile("data/orders.json")
	if err != nil {
		fmt.Println("Не удалось загрузить заказы:", err)
	}
	err = r.LoadUsersFromFile("data/users.json")
	if err != nil {
		fmt.Println("Не удалось загрузить пользователей:", err)
	}
	err = r.LoadDeliveriesFromFile("data/deliveries.json")
	if err != nil {
		fmt.Println("Не удалось загрузить доставки:", err)
	}
	err = r.LoadWarehousesFromFile("data/warehouses.json")
	if err != nil {
		fmt.Println("Не удалось загрузить склады:", err)
	}
	fmt.Println("Данные успешно загружены")
}

// метод, который ищет заказ по id

func (r *MemoryRepo) GetOrderByID(id string) (*model.Order, error) {
	r.muOrders.Lock()
	defer r.muOrders.Unlock()
	for _, order := range r.orders {
		if order.Id == id {
			return order, nil
		}
	}
	return nil, nil
}

// методы обновления статуса заказа

func (r *MemoryRepo) ConfirmOrder(orderId string) (bool, error) {
	r.muOrders.Lock()
	defer r.muOrders.Unlock()

	for _, order := range r.orders {
		if order.Id == orderId && order.Status == model.OrderCreated {
			order.Status = model.OrderConfirmed
			return true, nil
		}
	}
	return false, nil
}

func (r *MemoryRepo) DeliverOrder(orderId string) (bool, error) {
	r.muOrders.Lock()
	defer r.muOrders.Unlock()

	for _, order := range r.orders {
		if order.Id == orderId && order.Status == model.OrderConfirmed {
			order.Status = model.OrderDelivered
			return true, nil
		}
	}
	return false, nil
}

func (r *MemoryRepo) CancelOrder(orderId string) (bool, error) {
	r.muOrders.Lock()
	defer r.muOrders.Unlock()

	for _, order := range r.orders {
		if order.Id == orderId && (order.Status == model.OrderCreated || order.Status == model.OrderConfirmed) {
			order.Status = model.OrderCancelled
			return true, nil
		}
	}
	return false, nil
}

// метод удаления заказа

func (r *MemoryRepo) DeleteOrder(orderId string) (bool, error) {
	r.muOrders.Lock()
	for i, order := range r.orders {
		if order.Id == orderId {
			r.orders = append(r.orders[:i], r.orders[i+1:]...)
			r.muOrders.Unlock()
			if err := r.SaveOrdersToFile("data/orders.json"); err != nil {
				fmt.Println("Ошибка при сохранении заказов:", err)
			}
			return true, nil
		}
	}
	return false, nil
}

func (r *MemoryRepo) GetUserByID(id string) (*model.User, error) {
	r.muUsers.Lock()
	defer r.muUsers.Unlock()
	for _, user := range r.users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, nil
}

func (r *MemoryRepo) UpdateUserName(id, name string) (bool, error) {
	r.muUsers.Lock()
	defer r.muUsers.Unlock()

	for _, user := range r.users {
		if user.Id == id {
			user.Name = name
			return true, nil
		}
	}
	return false, nil // пользователь не найден
}

func (r *MemoryRepo) DeleteUser(id string) (bool, error) {
	r.muUsers.Lock()
	for i, user := range r.users {
		if user.Id == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			r.muUsers.Unlock()
			if err := r.SaveUsersToFile("data/users.json"); err != nil {
				fmt.Println("Ошибка при сохранении пользователей:", err)
			}
			return true, nil
		}
	}
	return false, nil
}

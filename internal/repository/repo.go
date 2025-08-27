package repository

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"order-ms/internal/model"
	"strconv"
	"time"
)

// Repo реализует интерфейс Repository
type Repo struct{}

// Создаем "конструктор"
func NewRepository() *Repo {
	return &Repo{}
}

// LogEvent сохраняет событие с TTL в Redis
func LogEvent(key, value string, ttl time.Duration) error {
	err := RedisClient.Set(Ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("не удалось записать в Redis: %w", err)
	}
	return nil
}

func (r *Repo) Save(s model.Storable) error {
	switch v := s.(type) {
	case *model.Order:
		return r.SaveOrder(v)
	case *model.User:
		return r.SaveUser(v)
	case *model.Delivery:
		return r.SaveDelivery(v)
	case *model.Warehouse:
		return r.SaveWarehouse(v)
	default:
		return fmt.Errorf("unsupported type")
	}
}

// Сохраняем новый заказ в MongoDB
func (r *Repo) SaveOrder(order *model.Order) error {
	_, err := OrderCollection.InsertOne(Ctx, order)
	if err != nil {
		return fmt.Errorf("не удалось сохранить заказ: %w", err)
	}

	// логируем событие в Redis
	key := fmt.Sprintf("order:%s:status", order.Id)
	value := strconv.Itoa(int(order.Status)) // конвертируем OrderStatus в строку числа
	if err := LogEvent(key, value, 24*time.Hour); err != nil {
		fmt.Println("Ошибка логирования создания нового заказа в Redis:", err)
	}

	return nil
}

// получаем все заказы из MongoDB
func (r *Repo) GetOrders() ([]*model.Order, error) {
	cursor, err := OrderCollection.Find(Ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(Ctx)

	var orders []*model.Order
	for cursor.Next(Ctx) {
		var order model.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

// получаем заказ по ID из MongoDB
func (r *Repo) GetOrderByID(id string) (*model.Order, error) {
	var order model.Order
	err := OrderCollection.FindOne(Ctx, bson.M{"id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // заказ не найден
		}
		return nil, fmt.Errorf("не удалось получить заказ: %w", err)
	}
	return &order, nil
}

// подтверждаем заказ в MongoDB
func (r *Repo) ConfirmOrder(orderId string) (bool, error) {
	filter := bson.M{"id": orderId, "status": model.OrderCreated} // ищем только созданный заказ
	update := bson.M{"$set": bson.M{"status": model.OrderConfirmed}}

	result, err := OrderCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return false, fmt.Errorf("не удалось подтвердить заказ: %w", err)
	}

	if result.MatchedCount == 0 {
		// либо заказа нет, либо статус не Created
		return false, nil
	}

	// логируем событие в Redis с TTL
	key := fmt.Sprintf("order:%s:status", orderId)
	value := strconv.Itoa(int(model.OrderConfirmed)) // конвертируем в строку
	if err := RedisClient.Set(Ctx, key, value, 24*time.Hour).Err(); err != nil {
		fmt.Println("Ошибка логирования подтверждения заказа в Redis:", err)
	}

	return true, nil
}

// отмечаем заказ как доставленный в MongoDB
func (r *Repo) DeliverOrder(orderId string) (bool, error) {
	filter := bson.M{"id": orderId, "status": model.OrderConfirmed} // ищем только подтверждённый заказ
	update := bson.M{"$set": bson.M{"status": model.OrderDelivered}}

	result, err := OrderCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return false, fmt.Errorf("не удалось отметить заказ как доставленный: %w", err)
	}

	if result.MatchedCount == 0 {
		// либо заказа нет, либо статус не Confirmed
		return false, nil
	}

	// логируем событие в Redis
	key := fmt.Sprintf("order:%s:status", orderId)
	value := strconv.Itoa(int(model.OrderDelivered))
	if err := RedisClient.Set(Ctx, key, value, 24*time.Hour).Err(); err != nil {
		fmt.Println("Ошибка логирования заказа в статусе доставлен в Redis:", err)
	}

	return true, nil
}

// отменяем заказ в MongoDB
func (r *Repo) CancelOrder(orderId string) (bool, error) {
	filter := bson.M{
		"id":     orderId,
		"status": bson.M{"$in": []int{int(model.OrderCreated), int(model.OrderConfirmed)}}, // можно отменять только созданные или подтверждённые
	}
	update := bson.M{"$set": bson.M{"status": model.OrderCancelled}}

	result, err := OrderCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return false, fmt.Errorf("не удалось отменить заказ: %w", err)
	}

	if result.MatchedCount == 0 {
		return false, nil
	}

	// логируем событие в Redis
	key := fmt.Sprintf("order:%s:status", orderId)
	value := strconv.Itoa(int(model.OrderCancelled))
	if err := RedisClient.Set(Ctx, key, value, 24*time.Hour).Err(); err != nil {
		fmt.Println("Ошибка логирования отмены заказа в Redis:", err)
	}

	return true, nil
}

// удаляем заказ в MongoDB
func (r *Repo) DeleteOrder(orderId string) (bool, error) {
	res, err := OrderCollection.DeleteOne(Ctx, bson.M{"id": orderId})
	if err != nil {
		return false, fmt.Errorf("ошибка при удалении заказа: %w", err)
	}
	if res.DeletedCount == 0 {
		return false, nil
	}
	key := fmt.Sprintf("order:%s:deleted", orderId)
	if err := RedisClient.Set(Ctx, key, "true", 24*time.Hour).Err(); err != nil {
		fmt.Println("Ошибка логирования удаления заказа в Redis:", err)
	}
	return true, nil
}

// сохраняем нового пользователя в MongoDB
func (r *Repo) SaveUser(user *model.User) error {
	_, err := UserCollection.InsertOne(Ctx, user)
	if err != nil {
		return fmt.Errorf("не удалось сохранить пользователя: %w", err)
	}

	// логируем событие в Redis
	key := fmt.Sprintf("user:%s:created", user.Id)
	if err := RedisClient.Set(Ctx, key, "true", 24*time.Hour).Err(); err != nil {
		fmt.Println("Ошибка логирования создания пользователя в Redis:", err)
	}

	return nil
}

// получаем всех пользователей
func (r *Repo) GetUsers() ([]*model.User, error) {
	cursor, err := UserCollection.Find(Ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(Ctx)

	var users []*model.User
	for cursor.Next(Ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// получаем пользователя по ID из MongoDB
func (r *Repo) GetUserByID(id string) (*model.User, error) {
	var user model.User
	err := UserCollection.FindOne(Ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // заказ не найден
		}
		return nil, fmt.Errorf("не удалось получить пользователя: %w", err)
	}
	return &user, nil
}

// обновляем имя пользователя
func (r *Repo) UpdateUserName(id, name string) (bool, error) {
	filter := bson.M{"id": id}                     // ищем пользователя по id
	update := bson.M{"$set": bson.M{"name": name}} // обновляем поле name
	result, err := UserCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return false, fmt.Errorf("не удалось обновить имя пользователя: %w", err)
	}

	if result.MatchedCount == 0 {
		// пользователь не найден
		return false, nil
	}

	// логируем изменение в Redis
	key := fmt.Sprintf("user:%s:name", id)
	if err := RedisClient.Set(Ctx, key, name, 24*time.Hour).Err(); err != nil {
		fmt.Println("Ошибка логирования в Redis:", err)
	}

	return true, nil
}

// удаляем пользователя в MongoDB
func (r *Repo) DeleteUser(id string) (bool, error) {
	res, err := UserCollection.DeleteOne(Ctx, bson.M{"id": id})
	if err != nil {
		return false, fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}
	if res.DeletedCount == 0 {
		return false, nil
	}
	key := fmt.Sprintf("user:%s:deleted", id)
	if err := RedisClient.Set(Ctx, key, "true", 24*time.Hour).Err(); err != nil {
		fmt.Println("Ошибка логирования в Redis:", err)
	}
	return true, nil
}

// сохраняем новую доставку в MongoDB
func (r *Repo) SaveDelivery(delivery *model.Delivery) error {
	_, err := DeliveryCollection.InsertOne(Ctx, delivery)
	if err != nil {
		return fmt.Errorf("не удалось сохранить доставку: %w", err)
	}
	return nil
}

// получаем все доставки
func (r *Repo) GetDeliveries() ([]*model.Delivery, error) {
	cursor, err := DeliveryCollection.Find(Ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(Ctx)

	var deliveries []*model.Delivery
	for cursor.Next(Ctx) {
		var delivery model.Delivery
		if err := cursor.Decode(&delivery); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, &delivery)
	}

	return deliveries, nil
}

// Сохраняем новый склад в MongoDB
func (r *Repo) SaveWarehouse(warehouse *model.Warehouse) error {
	_, err := WarehouseCollection.InsertOne(Ctx, warehouse)
	if err != nil {
		return fmt.Errorf("не удалось сохранить склад: %w", err)
	}
	return nil
}

// получаем все склады
func (r *Repo) GetWarehouses() ([]*model.Warehouse, error) {
	cursor, err := WarehouseCollection.Find(Ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(Ctx)

	var warehouses []*model.Warehouse
	for cursor.Next(Ctx) {
		var warehouse model.Warehouse
		if err := cursor.Decode(&warehouse); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &warehouse)
	}

	return warehouses, nil
}

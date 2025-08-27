package service

import (
	"context"
	"fmt"
	"order-ms/internal/model"
	"time"
)

// Service — обертка вокруг репозитория
type Service struct {
	repo Repository
}

// Создаем конструктор
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Logger выводит информацию о текущем состоянии базы
func (s *Service) Logger(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			orders, _ := s.repo.GetOrders()
			fmt.Printf("Orders in DB: %d\n", len(orders))

			users, _ := s.repo.GetUsers()
			fmt.Printf("Users in DB: %d\n", len(users))

			deliveries, _ := s.repo.GetDeliveries()
			fmt.Printf("Deliveries in DB: %d\n", len(deliveries))

			warehouses, _ := s.repo.GetWarehouses()
			fmt.Printf("Warehouses in DB: %d\n", len(warehouses))
		}
	}
}

// Save сохраняет объект через репозиторий
func (s *Service) Save(sObj model.Storable) error {
	return s.repo.Save(sObj)
}

//// MongoRepo реализует Repository для MongoDB и Redis
//type MongoRepo struct{}
//
//// Save сохраняет объект в MongoDB и логирует в Redis
//func (r *MongoRepo) Save(s model.Storable) error {
//	switch v := s.(type) {
//	case *model.Order:
//		if err := repository.SaveOrder(v); err != nil {
//			return err
//		}
//		// Логируем в Redis с TTL 1 час
//		key := fmt.Sprintf("order:%s", v.Id)
//		repository.RedisClient.Set(repository.Ctx, key, fmt.Sprintf("%+v", v), time.Hour)
//
//	case *model.User:
//		if err := repository.SaveUser(v); err != nil {
//			return err
//		}
//		key := fmt.Sprintf("user:%s", v.Id)
//		repository.RedisClient.Set(repository.Ctx, key, fmt.Sprintf("%+v", v), time.Hour)
//
//	case *model.Delivery:
//		if err := repository.SaveDelivery(v); err != nil {
//			return err
//		}
//		key := fmt.Sprintf("delivery:%d", v.Id)
//		repository.RedisClient.Set(repository.Ctx, key, fmt.Sprintf("%+v", v), time.Hour)
//
//	case *model.Warehouse:
//		if err := repository.SaveWarehouse(v); err != nil {
//			return err
//		}
//		key := fmt.Sprintf("warehouse:%d", v.Id)
//		repository.RedisClient.Set(repository.Ctx, key, fmt.Sprintf("%+v", v), time.Hour)
//
//	default:
//		return fmt.Errorf("неизвестный тип объекта: %T", s)
//	}
//
//	return nil
//}
//
//// Logger выводит информацию о текущем состоянии базы
//func Logger(ctx context.Context) {
//	ticker := time.NewTicker(5 * time.Second)
//	defer ticker.Stop()
//
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case <-ticker.C:
//			orders, _ := repository.GetOrders()
//			fmt.Printf("Orders in DB: %d\n", len(orders))
//
//			users, _ := repository.GetUsers()
//			fmt.Printf("Users in DB: %d\n", len(users))
//
//			deliveries, _ := repository.GetDeliveries()
//			fmt.Printf("Deliveries in DB: %d\n", len(deliveries))
//
//			warehouses, _ := repository.GetWarehouses()
//			fmt.Printf("Warehouses in DB: %d\n", len(warehouses))
//		}
//	}
//}

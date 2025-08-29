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

// NewService создаёт новый экземпляр Service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Logger выводит информацию о текущем состоянии базы
func (s *Service) Logger(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
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

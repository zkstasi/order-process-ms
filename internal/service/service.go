package service

import (
	"order-ms/internal/model"
	"order-ms/internal/repository"
)

// функция для создания структур и передачи их в репозиторий

func CreateStructs() {

	order := model.NewOrder("user123")
	user := model.NewUser("user123", "Петя")
	delivery := model.NewDelivery(65, "order-783", "user123", "ул. Ленина", 0)
	warehouse := model.NewWarehouse(543, "order-783", 0)

	repository.SaveStorable(order)
	repository.SaveStorable(user)
	repository.SaveStorable(delivery)
	repository.SaveStorable(warehouse)
}

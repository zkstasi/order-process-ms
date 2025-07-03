package service

import (
	"order-ms/internal/model"
	"order-ms/internal/repository"
)

// функция для создания структур и передачи их в репозиторий

func CreateStructs() {

	order := model.NewOrder("user123")
	user := model.NewUser(234, "Петя")
	delivery := model.NewDelivery(65, 78, 123, "ул. Ленина")
	warehouse := model.NewWarehouse(543, 78, 0)

	repository.SaveStorable(order)
	repository.SaveStorable(user)
	repository.SaveStorable(delivery)
	repository.SaveStorable(warehouse)
}

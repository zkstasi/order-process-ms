package service

import (
	"order-ms/internal/model"
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

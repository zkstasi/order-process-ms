package main

import (
	"order-ms/internal/model"
	"order-ms/internal/repository"
	"order-ms/internal/service"
	"sync"
	"time"
)

func main() {

	stop := make(chan struct{})
	dataChan := make(chan model.Storable)

	var wgCrSt sync.WaitGroup
	var wgSaSt sync.WaitGroup

	wgCrSt.Add(1) //запуск создателя структур
	go func() {
		defer wgCrSt.Done()
		service.CreateStructs(dataChan, stop)
	}()

	wgSaSt.Add(1) // запуск хранителя в репозиторий
	go func() {
		defer wgSaSt.Done()
		repository.SaveStorable(dataChan)
	}()

	time.Sleep(3 * time.Second) // работа 3 секунды

	close(stop)
	wgCrSt.Wait() // ждем завершения CreateStructs

	close(dataChan)
	wgSaSt.Wait() // ждем завершения SaveStorable
}

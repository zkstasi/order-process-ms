package main

import (
	"order-ms/internal/model"
	"order-ms/internal/repository"
	"order-ms/internal/service"
	"sync"
	"time"
)

func main() {

	loggerStop := make(chan struct{})
	stop := make(chan struct{})
	dataChan := make(chan model.Storable)

	var wgLog sync.WaitGroup
	var wgCrSt sync.WaitGroup
	var wgSaSt sync.WaitGroup

	wgLog.Add(1) // запуск логирования
	go func() {
		defer wgLog.Done()
		service.Logger(loggerStop)
	}()

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

	close(stop)   // останавливаем создателя структур
	wgCrSt.Wait() // ждем завершения CreateStructs

	close(dataChan) // закрываем канал для SaveStorable
	wgSaSt.Wait()   // ждем завершения SaveStorable

	close(loggerStop) // останавливаем Logger
	wgLog.Wait()      // Ждем завершения Logger

}

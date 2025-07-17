package main

import (
	"context"
	"order-ms/internal/model"
	"order-ms/internal/repository"
	"order-ms/internal/service"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	repository.LoadAllData() // загружаем все из файлов при запуске

	//создание контекста, который отменится, когда пользователь нажмет Ctrl+C или придет другой сигнал завершения
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // освобождаем ресурсы

	dataChan := make(chan model.Storable)

	var wgLog sync.WaitGroup
	var wgCrSt sync.WaitGroup
	var wgSaSt sync.WaitGroup

	wgLog.Add(1) // запуск логирования
	go func() {
		defer wgLog.Done()
		service.Logger(ctx)
	}()

	wgCrSt.Add(1) //запуск создателя структур
	go func() {
		defer wgCrSt.Done()
		service.CreateStructs(ctx, dataChan)
	}()

	wgSaSt.Add(1) // запуск хранителя в репозиторий
	go func() {
		defer wgSaSt.Done()
		service.ProcessDataChan(dataChan)
	}()

	<-ctx.Done() // ждем сигнала ОС

	wgCrSt.Wait() // ждем завершения CreateStructs

	close(dataChan) // закрываем канал для ProcessDataChan
	wgSaSt.Wait()   // ждем завершения ProcessDataChan

	wgLog.Wait() // Ждем завершения Logger

	repository.SaveAllData()
}

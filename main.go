// @title Order Processing API
// @version 1.0
// @description API для управления пользователями и заказами
// @host localhost:8080
// @BasePath /

package main

import (
	"context"
	"log"
	_ "order-ms/docs"
	"order-ms/internal/repository"
	"order-ms/internal/service"
	"order-ms/internal/web"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	repository.LoadAllData() // загружаем все из файлов при запуске

	//создание контекста, который отменится, когда пользователь нажмет Ctrl+C или придет другой сигнал завершения
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // освобождаем ресурсы

	var wgLog sync.WaitGroup

	wgLog.Add(1) // запуск логирования
	go func() {
		defer wgLog.Done()
		service.Logger(ctx)
	}()

	// запуск http-сервера
	webServer := web.NewServer(":8080")
	go func() {
		if err := webServer.Start(); err != nil {
			log.Printf("Server start error: %v\n", err)
		}
	}()

	<-ctx.Done() // ждем сигнала ОС
	wgLog.Wait() // Ждем завершения Logger

	repository.SaveAllData()
}

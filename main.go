// @title Order Processing API
// @version 1.0
// @description API для управления пользователями и заказами
// @host localhost:8080
// @BasePath /

package main

import (
	"context"
	"flag"
	"log"
	_ "order-ms/docs"
	"order-ms/internal/repository/memory"
	repository "order-ms/internal/repository/nosql"
	"order-ms/internal/service"
	"order-ms/internal/web"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Флаг командной строки для выбора репозитория
	useMemory := flag.Bool("memory", false, "Use in-memory repository")
	flag.Parse()

	//создание контекста, который отменится, когда пользователь нажмет Ctrl+C или придет другой сигнал завершения
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // освобождаем ресурсы

	var repo service.Repository
	if *useMemory {
		repo = memory.NewMemoryRepo()
	} else {
		// Инициализация Mongo/Redis для настоящего репозитория
		if err := repository.InitDB(); err != nil {
			log.Fatalf("Ошибка инициализации базы данных: %v", err)
		}
		defer repository.CloseDB()
		repo = repository.NewRepository()
	}

	// Создаем сервис с выбранным репозиторием
	svc := service.NewService(repo)

	var wg sync.WaitGroup

	wg.Add(1) // запуск логирования
	go func() {
		defer wg.Done()
		svc.Logger(ctx)
	}()

	// запуск http-сервера
	webServer := web.NewServer(":8080", repo)
	go func() {
		if err := webServer.Start(); err != nil {
			log.Printf("Server start error: %v\n", err)
		}
	}()

	//// Запускаем gRPC сервер
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//
	//	lis, err := net.Listen("tcp", ":50051") // порт для gRPC, можно другой
	//	if err != nil {
	//		log.Fatalf("Failed to listen: %v", err)
	//	}
	//
	//	grpcServer := grpcServerPkg.NewGrpcServer() // твой сервер
	//
	//	go func() {
	//		<-ctx.Done() // на сигнал завершения
	//		log.Println("Stopping gRPC server...")
	//		grpcServer.GracefulStop()
	//	}()
	//
	//	log.Println("Starting gRPC server on :50051")
	//	if err := grpcServer.Serve(lis); err != nil {
	//		log.Printf("gRPC server error: %v", err)
	//	}
	//}()

	//// Запускаем тест клиента (после небольшого ожидания)
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//
	//	time.Sleep(2 * time.Second) // подождать, пока сервер запустится
	//
	//	grpcClient, err := service.NewGrpcClient("localhost:50051")
	//	if err != nil {
	//		log.Printf("Failed to create gRPC client: %v", err)
	//		return
	//	}
	//	defer grpcClient.Close()
	//
	//	// Создаем пользователя
	//	user, err := grpcClient.CreateUserExample()
	//	if err != nil {
	//		log.Printf("CreateUserExample error: %v", err)
	//		return
	//	}
	//	log.Printf("Created user: %v", user)
	//
	//	// Получаем пользователя по ID
	//	gotUser, err := grpcClient.GetUserExample(user.Id)
	//	if err != nil {
	//		log.Printf("GetUserExample error: %v", err)
	//	} else {
	//		log.Printf("Got user: %v", gotUser)
	//	}
	//
	//	// Обновляем имя пользователя
	//	updatedUser, err := grpcClient.UpdateUserExample(user.Id, "Bob")
	//	if err != nil {
	//		log.Printf("UpdateUserExample error: %v", err)
	//	} else {
	//		log.Printf("Updated user: %v", updatedUser)
	//	}
	//
	//	// Список всех пользователей
	//	users, err := grpcClient.ListUsersExample()
	//	if err != nil {
	//		log.Printf("ListUsersExample error: %v", err)
	//	} else {
	//		log.Printf("Users list: %v", users)
	//	}
	//
	//	// Создаем заказ для пользователя
	//	order, err := grpcClient.CreateOrderExample(user.Id)
	//	if err != nil {
	//		log.Printf("CreateOrderExample error: %v", err)
	//	} else {
	//		log.Printf("Created order: %v", order)
	//	}
	//
	//	// Подтверждаем заказ
	//	confirmedOrder, err := grpcClient.ConfirmOrderExample(order.Id)
	//	if err != nil {
	//		log.Printf("ConfirmOrderExample error: %v", err)
	//	} else {
	//		log.Printf("Confirmed order: %v", confirmedOrder)
	//	}
	//
	//	// Отправляем заказ на доставку
	//	deliveredOrder, err := grpcClient.DeliverOrderExample(order.Id)
	//	if err != nil {
	//		log.Printf("DeliverOrderExample error: %v", err)
	//	} else {
	//		log.Printf("Delivered order: %v", deliveredOrder)
	//	}
	//
	//	// Отменяем заказ (если нужно)
	//	err = grpcClient.CancelOrderExample(order.Id)
	//	if err != nil {
	//		log.Printf("CancelOrderExample error: %v", err)
	//	} else {
	//		log.Printf("Order cancelled successfully")
	//	}
	//
	//	// Удаляем пользователя
	//	err = grpcClient.DeleteUserExample(user.Id)
	//	if err != nil {
	//		log.Printf("DeleteUserExample error: %v", err)
	//	} else {
	//		log.Printf("User deleted successfully")
	//	}
	//}()

	<-ctx.Done() // ждем сигнала ОС
	wg.Wait()    // Ждем завершения горутин

	// Сохраняем данные MemoryRepo при завершении
	if memRepo, ok := repo.(*memory.MemoryRepo); ok {
		memRepo.SaveAllData()
	}

	log.Println("Приложение завершено")
}

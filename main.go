package main

import (
	"fmt"
	"order-ms/internal/model"
	"time"
)

// Для проверки
func main() {

	order := model.NewOrder("user123")

	// Выводим информацию о заказе для проверки
	fmt.Printf("Создан новый заказ:\n")
	fmt.Printf("ID: %s\n", order.Id())
	fmt.Printf("UserID: %s\n", order.UserId())
	fmt.Printf("Status: %d\n", order.Status())
	fmt.Printf("CreatedAt: %s\n", order.CreatedAt().Format(time.RFC3339))

	user := model.NewUser(456, "John")

	// Выводим информацию о пользователе для проверки
	fmt.Printf("Создан новый пользователь:\n")
	fmt.Printf("UserId: %d\n", user.Id())
	fmt.Printf("UserName: %s\n", user.Name())
}

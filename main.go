package main

import (
	"fmt"
	"order-ms/internal/model/order"
	"order-ms/internal/model/user"
)

// Для проверки
func main() {

	myOrder := order.Order{}
	myUser := user.User{}

	myOrder.SetID(123)
	myOrder.SetStatus(1)

	myUser.SetID(15)
	myUser.SetName("Настя")

	fmt.Printf("ID: %d, Status: %d\n", myOrder.GetID(), myOrder.GetStatus())
	fmt.Printf("ID: %d, Name: %s\n", myUser.GetID(), myUser.GetName())
}

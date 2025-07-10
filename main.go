package main

import (
	"fmt"
	"order-ms/internal/model"
	"order-ms/internal/repository"
	"order-ms/internal/service"
)

func main() {

	stop := make(chan struct{})
	dataChan := make(chan model.Storable)

	go service.CreateStructs(dataChan, stop)
	close(stop)

	go repository.SaveStorable(dataChan)

	fmt.Println("Starting service...")

}

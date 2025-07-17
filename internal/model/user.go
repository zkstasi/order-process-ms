package model

import (
	"fmt"
	"time"
)

type User struct {
	Id   string `json:"id"`   // Уникальный номер пользователя
	Name string `json:"name"` // Имя пользователя
}

// NewUser создаёт нового пользователя с заданным id и именем.

func NewUser(newName string) *User {
	return &User{
		Id:   generateUserID(),
		Name: newName,
	}
}

func generateUserID() string {
	return fmt.Sprintf("User-%d", time.Now().UnixNano())
}

// реализация интерфейса Storable

func (u *User) GetType() string {
	return "user"
}

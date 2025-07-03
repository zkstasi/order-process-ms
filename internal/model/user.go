package model

type User struct {
	id   string // Приватный уникальный номер пользователя
	name string // Имя пользователя
}

// NewUser создаёт нового пользователя с заданным id и именем.

func NewUser(newID string, newName string) *User {
	return &User{
		id:   newID,
		name: newName,
	}
}

func (u *User) Id() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) SetName(newName string) {
	u.name = newName
}

// реализация интерфейса Storable

func (u *User) GetType() string {
	return "user"
}

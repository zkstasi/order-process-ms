package model

type User struct {
	id   int    // Приватный уникальный номер пользователя
	name string // Имя пользователя
}

// NewUser создаёт нового пользователя с заданным id и именем.

func NewUser(newID int, newName string) *User {
	return &User{
		id:   newID,
		name: newName,
	}
}

func (u *User) Id() int {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) SetName(newName string) {
	u.name = newName
}

func (u *User) GetType() string {
	return "user"
}

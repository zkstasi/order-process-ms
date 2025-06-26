package model

type User struct {
	id   int    // Приватный уникальный номер пользователя
	name string // Имя пользователя
}

func (u *User) ID() int {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Fields() (int, string) {
	return u.id, u.name
}

func (u *User) SetName(newName string) {
	u.name = newName
}

func NewUser(newID int, newName string) *User {
	return &User{
		id:   newID,
		name: newName,
	}
}

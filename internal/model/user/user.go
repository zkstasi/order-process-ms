package user

type User struct {
	id   int    // Приватный уникальный номер
	name string // Имя пользователя
}

func (u *User) GetID() int {
	return u.id
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetAll() (int, string) {
	return u.id, u.name
}

func (u *User) SetID(newID int) {
	u.id = newID
}

func (u *User) SetName(newName string) {
	u.name = newName
}

func (u *User) SetAll(newID int, newName string) {
	u.id = newID
	u.name = newName
}

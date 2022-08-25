package models

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Active   bool   `db:"active"`
}

type UserService interface {
	Create(user *User) (User, error)
	FindByID(id int) (User, error)
	FindByUsername(username string) (User, error)
}

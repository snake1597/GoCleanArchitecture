package entities

import "time"

type User struct {
	ID           int       `gorm:"primaryKey; id" json:"id"`
	Account      string    `gorm:"account; unique; type:varchar(100); not null;" json:"account" `
	Password     string    `gorm:"password; type:varchar(30); not null;" json:"password" `
	Name         string    `gorm:"name; type:varchar(50); " json:"name" `
	Birthday     string    `gorm:"birthday; type:varchar(10);" json:"birthday" `
	RefreshToken string    `gorm:"refresh_token; type:varchar(30);" json:"refresh_token" `
	CreatedAt    time.Time `gorm:"created_at" json:"created_at" `
	UpdatedAt    time.Time `gorm:"updated_at" json:"updated_at" `
}

type UserRepository interface {
	Register(request User) (err error)
	Login(request User) (user User, err error)
	GetUser(id string) (user User, err error)
	GetAllUser() (userList []User, err error)
	UpdateUser(id string, data map[string]interface{}) (err error)
	DeleteUser(id string) (err error)
}

type UserUsecase interface {
	Register(user User) (err error)
	Login(user User) (id string, err error)
	GetUser(id string) (user User, err error)
	GetAllUser() (userList []User, err error)
	UpdateUser(id string, user User) (err error)
	DeleteUser(id string) (err error)
}

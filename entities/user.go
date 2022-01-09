package entities

import "time"

type User struct {
	ID           int       `gorm:"primaryKey; id" json:"id"`
	Account      string    `gorm:"account; unique; type:varchar(100); not null;" json:"account" `
	Password     string    `gorm:"password; type:varchar(30); not null;" json:"password" `
	FirstName    string    `gorm:"first_name; type:varchar(50); not null;" json:"first_name" `
	LastName     string    `gorm:"last_name; type:varchar(50); not null;" json:"last_name" `
	Birthday     time.Time `gorm:"birthday; not null;" json:"birthday" `
	RefreshToken string    `gorm:"refresh_token; type:varchar(30);" json:"refresh_token" `
	CreatedAt    time.Time `gorm:"created_at;" json:"created_at" `
	UpdatedAt    time.Time `gorm:"updated_at;" json:"updated_at" `
}

type UserRepository interface {
	Register(request *User) (err error)
	Login(request *User) (user *User, err error)
	GetUser(id string) (user *User, err error)
	GetAllUser() (userList []User, err error)
	UpdateUser(id string, data map[string]interface{}) (err error)
	DeleteUser(id string) (err error)
}

type UserUsecase interface {
	Register(request *User) (err error)
	Login(request *User) (id string, err error)
	GetUser(id string) (user *User, err error)
	GetAllUser() (userList []User, err error)
	UpdateUser(id string, user *User) (err error)
	DeleteUser(id string) (err error)
}

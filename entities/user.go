package entities

import "time"

// swagger:model
type User struct {
	// min: 1
	ID int `gorm:"primaryKey; id" json:"id"`

	// max length: 100
	// pattern: ^([a-zA-Z0-9.])*@([a-zA-Z0-9])*\.([a-zA-Z0-9])*
	Account string `gorm:"account; unique; type:varchar(100); not null;" json:"account" `

	// min length: 8
	// max length: 30
	// pattern: ^([a-zA-Z0-9]){8,30}$
	Password string `gorm:"password; type:varchar(30); not null;" json:"password" `

	// max length: 50
	FirstName string `gorm:"first_name; type:varchar(50); not null;" json:"first_name" `

	// max length: 50
	LastName string `gorm:"last_name; type:varchar(50); not null;" json:"last_name" `

	Birthday time.Time `gorm:"birthday; not null;" json:"birthday" `

	// max length: 30
	RefreshToken string `gorm:"refresh_token; type:varchar(30);" json:"refresh_token" `

	CreatedAt time.Time `gorm:"created_at;" json:"created_at" `

	UpdatedAt time.Time `gorm:"updated_at;" json:"updated_at" `
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

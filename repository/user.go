package repository

import (
	"GoCleanArchitecture/entities"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) entities.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Register(user *entities.User) (err error) {
	err = r.db.Create(&user).Error
	return err
}

func (r *userRepository) Login(request *entities.User) (user *entities.User, err error) {
	err = r.db.Table("users").Where("account = ?", request.Account).First(&user).Error
	return user, err
}

func (r *userRepository) GetUser(id string) (user *entities.User, err error) {
	err = r.db.Table("users").Where("id = ?", id).Omit("password").Find(&user).Error
	return user, err
}

func (r *userRepository) GetAllUser() (userList []entities.User, err error) {
	err = r.db.Table("users").Find(&userList).Error
	return userList, err
}

func (r *userRepository) UpdateUser(id string, data map[string]interface{}) (err error) {
	err = r.db.Table("users").Where("id = ?", id).Updates(data).Error
	return err
}

func (r *userRepository) DeleteUser(id string) (err error) {
	err = r.db.Table("users").Where("id = ?", id).Delete(&entities.User{}).Error
	return err
}

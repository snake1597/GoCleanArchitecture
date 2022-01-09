package usecase

import (
	"GoCleanArchitecture/entities"
	"fmt"
	"regexp"
	"strconv"
)

type userUsecase struct {
	userRepo entities.UserRepository
}

var (
	regEmail    *regexp.Regexp
	regPassword *regexp.Regexp
)

func NewUserUsecase(userRepo entities.UserRepository) entities.UserUsecase {
	emailPattern := `^([a-zA-Z0-9.])*@([a-zA-Z0-9])*\.([a-zA-Z0-9])*`
	regEmail = regexp.MustCompile(emailPattern)
	passwordPattern := `^([a-zA-Z0-9]){8,30}$`
	regPassword = regexp.MustCompile(passwordPattern)

	return &userUsecase{userRepo}
}

func (u *userUsecase) Register(user *entities.User) (err error) {
	if user.Account == "" {
		return fmt.Errorf("account can not be empty")
	}

	if user.Password == "" {
		return fmt.Errorf("password can not be empty")
	}

	if user.FirstName == "" {
		return fmt.Errorf("first name can not be empty")
	}

	if user.LastName == "" {
		return fmt.Errorf("last name can not be empty")
	}

	if user.Birthday == "" {
		return fmt.Errorf("birthday can not be empty")
	}

	if !verifyEmailFormat(user.Account) {
		return fmt.Errorf("incorrect email format")
	}

	if !verifyPasswordFormat(user.Password) {
		return fmt.Errorf("password must be at 8 to 30")
	}

	err = u.userRepo.Register(user)
	return err
}

func (u *userUsecase) Login(request *entities.User) (userId string, err error) {
	user, err := u.userRepo.Login(request)
	if err != nil {
		return "", fmt.Errorf("account is not exist")
	}

	if request.Password != user.Password {
		return "", fmt.Errorf("password is not valid")
	}

	userId = strconv.Itoa(user.ID)
	return userId, err
}

func (u *userUsecase) GetUser(id string) (user *entities.User, err error) {
	user, err = u.userRepo.GetUser(id)
	if err != nil {
		return &entities.User{}, err
	}

	return user, nil
}

func (u *userUsecase) GetAllUser() (userList []entities.User, err error) {
	userList, err = u.userRepo.GetAllUser()
	if err != nil {
		return nil, err
	}

	return userList, err
}

func (u *userUsecase) UpdateUser(id string, user *entities.User) (err error) {
	data := make(map[string]interface{})

	if user.FirstName != "" {
		data["firstName"] = user.FirstName
	}
	if user.LastName != "" {
		data["lastName"] = user.LastName
	}
	if user.Birthday != "" {
		data["birthday"] = user.Birthday
	}

	err = u.userRepo.UpdateUser(id, data)
	if err != nil {
		return fmt.Errorf("update failed")
	}

	return nil
}

func (u *userUsecase) DeleteUser(id string) (err error) {
	err = u.userRepo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("delete failed")
	}

	return nil
}

func verifyEmailFormat(email string) bool {
	return regEmail.MatchString(email)
}

func verifyPasswordFormat(email string) bool {
	return regPassword.MatchString(email)
}

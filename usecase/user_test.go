package usecase_test

import (
	"fmt"
	"testing"

	"GoCleanArchitecture/entities"
	"GoCleanArchitecture/entities/mocks"
	"GoCleanArchitecture/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		mockUser := &entities.User{
			Account:  "a@gmail.com",
			Password: "12345678",
		}

		mockUserRepo.
			On("Register", mock.Anything).
			Return(nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		err := u.Register(mockUser)
		assert.NoError(t, err)
	})

	t.Run("email is null", func(t *testing.T) {
		mockUser := &entities.User{
			Account:  "",
			Password: "12345678",
		}

		u := usecase.NewUserUsecase(mockUserRepo)
		err := u.Register(mockUser)
		assert.Error(t, err)
	})

	t.Run("pw is null", func(t *testing.T) {
		mockUser := &entities.User{
			Account:  "a@gmail.com",
			Password: "",
		}

		u := usecase.NewUserUsecase(mockUserRepo)
		err := u.Register(mockUser)
		assert.Error(t, err)
	})

	t.Run("incorrect email format", func(t *testing.T) {
		mockUser := &entities.User{
			Account:  "a@gmail",
			Password: "12345678",
		}

		u := usecase.NewUserUsecase(mockUserRepo)
		err := u.Register(mockUser)
		assert.Error(t, err)
	})

	t.Run("incorrect pw format", func(t *testing.T) {
		mockUser := &entities.User{
			Account:  "a@gmail.com",
			Password: "123456",
		}

		u := usecase.NewUserUsecase(mockUserRepo)
		err := u.Register(mockUser)
		assert.Error(t, err)

		mockUser = &entities.User{
			Account:  "a@gmail.com",
			Password: "123456789123456789123456789123456789", // 36 numbers
		}
		err = u.Register(mockUser)
		assert.Error(t, err)
	})
}

func TestLogin(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		mockUser := &entities.User{
			Account:  "a@gmail.com",
			Password: "12345678",
		}

		mockUserRepo.
			On("Login", mock.Anything).
			Return(entities.User{Account: "a@gmail.com", Password: "12345678"}, nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		data, err := u.Login(mockUser)
		assert.NoError(t, err)
		assert.NotNil(t, data)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("invalid account", func(t *testing.T) {
		mockUser := &entities.User{
			Account:  "a@gmail.com",
			Password: "12345678",
		}

		mockUserRepo.
			On("Login", mock.Anything).
			Return(entities.User{}, fmt.Errorf("account is not exist")).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		data, err := u.Login(mockUser)

		assert.Error(t, err)
		assert.Equal(t, data, "")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("invalid pw", func(t *testing.T) {
		mockUser := &entities.User{
			Account:  "a@gmail.com",
			Password: "12345678",
		}

		mockUserRepo.
			On("Login", mock.Anything).
			Return(entities.User{Account: "a@gmail.com", Password: "12345"}, nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		data, err := u.Login(mockUser)

		assert.Error(t, err)
		assert.Equal(t, data, "")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := entities.User{
		ID: 9,
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("GetUser", mock.AnythingOfType("string")).
			Return(mockUser, nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		user, err := u.GetUser("9")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		mockUserRepo.
			On("GetUser", mock.AnythingOfType("string")).
			Return(entities.User{}, fmt.Errorf("password is not valid")).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		user, err := u.GetUser("9")

		assert.Error(t, err)
		assert.Empty(t, user)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGetAllUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := []entities.User{
		{
			ID:        8,
			FirstName: "firstUser",
		},
		{
			ID:        9,
			FirstName: "secondUser",
		},
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("GetAllUser", mock.AnythingOfType("string")).
			Return(mockUser, nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		user, err := u.GetAllUser()

		assert.NoError(t, err)
		assert.NotNil(t, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		mockUserRepo.
			On("GetAllUser", mock.AnythingOfType("string")).
			Return([]entities.User{}, fmt.Errorf("error")).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		user, err := u.GetAllUser()

		assert.Error(t, err)
		assert.Empty(t, user)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &entities.User{
		ID: 9,
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("UpdateUser", mock.AnythingOfType("string"), mock.Anything).
			Return(nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		err := u.UpdateUser("9", mockUser)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		mockUserRepo.
			On("UpdateUser", mock.AnythingOfType("string"), mock.Anything).
			Return(fmt.Errorf("update failed")).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		err := u.UpdateUser("9", mockUser)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.
			On("DeleteUser", mock.AnythingOfType("string")).
			Return(nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		err := u.DeleteUser("9")

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		mockUserRepo.
			On("DeleteUser", mock.AnythingOfType("string")).
			Return(fmt.Errorf("update failed")).Once()

		u := usecase.NewUserUsecase(mockUserRepo)
		err := u.DeleteUser("9")

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

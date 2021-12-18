package usecase_test

import (
	"GoCleanArchitecture/entities"
	"GoCleanArchitecture/entities/mocks"
	"GoCleanArchitecture/usecase"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateToken(t *testing.T) {
	mockTokenRepo := new(mocks.TokenRepository)

	t.Run("success", func(t *testing.T) {
		mockTokenRepo.
			On("UpdateRefreshToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(nil).Once()

		u := usecase.NewTokenUsecase("jwtKey", mockTokenRepo)
		token, err := u.CreateToken("9")
		t.Log(token.AccessToken)
		t.Log(token.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, token)
	})
}

func TestVerifyToken(t *testing.T) {
	mockTokenRepo := new(mocks.TokenRepository)

	t.Run("success", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := make(jwt.MapClaims)
		claims["id"] = "9"
		claims["iat"] = time.Now().Unix()
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		token.Claims = claims
		accessToken, _ := token.SignedString([]byte("jwtKey"))

		u := usecase.NewTokenUsecase("jwtKey", mockTokenRepo)
		data, err := u.VerifyToken(accessToken)

		assert.NoError(t, err)
		assert.NotNil(t, data)
	})

	t.Run("fail", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := make(jwt.MapClaims)
		claims["id"] = "9"
		claims["iat"] = time.Now().Unix()
		claims["exp"] = time.Now().Add(time.Minute * -15).Unix()
		token.Claims = claims
		accessToken, _ := token.SignedString([]byte("jwtKey"))

		u := usecase.NewTokenUsecase("jwtKey", mockTokenRepo)
		data, err := u.VerifyToken(accessToken)

		assert.Error(t, err)
		assert.Nil(t, data)
	})
}

func TestRefreshAccessToken(t *testing.T) {
	mockTokenRepo := new(mocks.TokenRepository)

	t.Run("success", func(t *testing.T) {
		mockTokenRepo.
			On("CheckRefreshToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(entities.User{ID: 9, Account: "a@gmail.com"}, nil).Once()

		mockTokenRepo.
			On("UpdateRefreshToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(nil).Once()

		u := usecase.NewTokenUsecase("jwtkey", mockTokenRepo)
		token, err := u.RefreshAccessToken("userId", "refreshToken")

		assert.NoError(t, err)
		assert.NotNil(t, token)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("refresh token is not exist", func(t *testing.T) {
		mockTokenRepo.
			On("CheckRefreshToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(entities.User{}, fmt.Errorf("refresh token is not exist")).Once()

		u := usecase.NewTokenUsecase("jwtkey", mockTokenRepo)
		token, err := u.RefreshAccessToken("userId", "refreshToken")

		assert.Error(t, err)
		assert.Empty(t, token)
		mockTokenRepo.AssertExpectations(t)
	})
}

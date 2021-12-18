package http_test

import (
	"GoCleanArchitecture/entities"
	"GoCleanArchitecture/entities/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	delivery "GoCleanArchitecture/delivery/http"
	"GoCleanArchitecture/delivery/http/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	mockUserUsecase := new(mocks.UserUsecase)
	mockTokenUsecase := new(mocks.TokenUsecase)
	authMiddleware := middlewares.NewAuthMiddlewares(mockTokenUsecase)

	mockToken := entities.Token{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}

	handler := &delivery.UserHandler{mockUserUsecase, mockTokenUsecase, authMiddleware}
	r := gin.Default()
	r.POST("/api/v1/users/login", handler.Login)

	t.Run("success", func(t *testing.T) {
		data := map[string]string{
			"user_id":       "userId",
			"access_token":  mockToken.AccessToken,
			"refresh_token": mockToken.RefreshToken,
		}
		response := map[string]interface{}{"status": "success", "data": data}
		mockResponse, _ := json.Marshal(response)

		mockUserUsecase.On("Login", mock.Anything).Return("userId", nil).Once()
		mockTokenUsecase.On("CreateToken", mock.AnythingOfType("string")).Return(mockToken, nil).Once()

		jsonString := []byte(`{"account":"arthur", "password":"12345678"}`)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(jsonString))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(mockResponse), w.Body.String())
	})

	t.Run("no data from the reruest", func(t *testing.T) {
		mockUserUsecase.On("Login", mock.Anything).Return("userId", nil).Once()
		mockTokenUsecase.On("CreateToken", mock.AnythingOfType("string")).Return(mockToken, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/login", nil)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetUser(t *testing.T) {
	mockUserUsecase := new(mocks.UserUsecase)
	mockTokenUsecase := new(mocks.TokenUsecase)
	authMiddleware := middlewares.NewAuthMiddlewares(mockTokenUsecase)

	mockUser := entities.User{
		ID: 9,
	}

	handler := &delivery.UserHandler{mockUserUsecase, mockTokenUsecase, authMiddleware}
	r := gin.Default()
	r.GET("/api/v1/users/:userId", handler.GetUser)

	t.Run("success", func(t *testing.T) {
		response := map[string]interface{}{"status": "success", "data": mockUser}
		mockResponse, _ := json.Marshal(response)
		mockUserUsecase.On("GetUser", mock.AnythingOfType("string")).Return(mockUser, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/userId", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(mockResponse), w.Body.String())
	})

	t.Run("fail", func(t *testing.T) {
		response := map[string]interface{}{"status": "failed", "error": "invalid id"}
		mockResponse, _ := json.Marshal(response)
		mockUserUsecase.On("GetUser", mock.AnythingOfType("string")).Return(entities.User{}, fmt.Errorf("invalid id")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/userId", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(mockResponse), w.Body.String())
	})
}

package http_test

import (
	delivery "GoCleanArchitecture/delivery/http"
	"GoCleanArchitecture/entities"
	"GoCleanArchitecture/entities/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRefreshAccessToken(t *testing.T) {
	mockTokenUsecase := new(mocks.TokenUsecase)

	mockToken := entities.Token{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}

	handler := &delivery.TokenHandler{mockTokenUsecase}
	r := gin.Default()
	r.PATCH("/api/v1/token/refresh/:userId", handler.RefreshAccessToken)

	t.Run("success", func(t *testing.T) {
		data := map[string]string{
			"user_id":       "userId",
			"access_token":  mockToken.AccessToken,
			"refresh_token": mockToken.RefreshToken,
		}
		response := map[string]interface{}{"status": "success", "data": data}
		mockResponse, _ := json.Marshal(response)

		mockTokenUsecase.On("RefreshAccessToken", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockToken, nil).Once()

		var bearer = "Bearer " + "refreshToken"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/token/refresh/userId", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", bearer)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(mockResponse), w.Body.String())
	})
}

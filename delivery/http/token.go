package http

import (
	"GoCleanArchitecture/entities"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	TokenUsecase entities.TokenUsecase
}

func NewTokenHandler(router *gin.Engine, tokenUsecase entities.TokenUsecase) {
	handler := &TokenHandler{tokenUsecase}

	r := router.Group("/api/v1/token")
	{
		r.PATCH("/refresh/:userId", handler.RefreshAccessToken)
	}
}

func (h *TokenHandler) RefreshAccessToken(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "token is missing"})
		return
	}

	refreshToken := strings.Split(auth, "Bearer ")[1]
	userId := c.Param("userId")

	token, err := h.TokenUsecase.RefreshAccessToken(userId, refreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	data := map[string]string{
		"user_id":       userId,
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

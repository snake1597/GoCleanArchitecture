package middlewares

import (
	"GoCleanArchitecture/docs/swagger"
	"GoCleanArchitecture/entities"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddlewares struct {
	tokenUsecase entities.TokenUsecase
}

func NewAuthMiddlewares(tokenUsecase entities.TokenUsecase) *AuthMiddlewares {
	return &AuthMiddlewares{tokenUsecase}
}

func (m *AuthMiddlewares) ParseToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, swagger.GenericError{Status: "failed", Message: "token is require"})
			return
		}

		token := strings.Split(auth, "Bearer ")[1]

		tokenInfo, err := m.tokenUsecase.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, swagger.GenericError{Status: "failed", Message: err.Error()})
			return
		}

		userId := c.Param("userId")
		if userId != "" && userId != tokenInfo["user_id"] {
			c.AbortWithStatusJSON(http.StatusUnauthorized, swagger.GenericError{Status: "failed", Message: "invalid user id"})
			return
		}

		c.Set("userId", tokenInfo["user_id"])
		c.Next()
	}
}

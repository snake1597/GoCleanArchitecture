package http

import (
	"GoCleanArchitecture/docs/swagger"
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
		r.PATCH("/token/:userId", handler.RefreshAccessToken)
	}
}

// swagger:operation PATCH /token/{userId} token refreshAccessToken
// ---
// summary: Refresh the access token when it was expired.
// parameters:
// - name: userId
//   in: path
//   description: user id
//   type: string
//   required: true
// security:
// - Bearer: []
// responses:
//   "200":
//     "$ref": "#/responses/loginResponse"
//   "400":
//     "$ref": "#/responses/genericError"
func (h *TokenHandler) RefreshAccessToken(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: "token is missing"})
		return
	}

	refreshToken := strings.Split(auth, "Bearer ")[1]
	userId := c.Param("userId")

	token, err := h.TokenUsecase.RefreshAccessToken(userId, refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: err.Error()})
		return
	}

	response := entities.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	c.JSON(http.StatusOK, swagger.LoginResponse{Status: "success", Data: response})
}

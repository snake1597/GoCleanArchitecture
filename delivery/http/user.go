package http

import (
	"GoCleanArchitecture/delivery/http/middlewares"
	"GoCleanArchitecture/docs/swagger"
	"GoCleanArchitecture/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUsecase    entities.UserUsecase
	TokenUsecase   entities.TokenUsecase
	AuthMiddleware *middlewares.AuthMiddlewares
}

func NewUserHandler(router *gin.Engine, userUsecase entities.UserUsecase, tokenUsecase entities.TokenUsecase, authMiddleware *middlewares.AuthMiddlewares) {
	handler := &UserHandler{userUsecase, tokenUsecase, authMiddleware}

	r := router.Group("/api/v1/users")
	{
		r.POST("/register", handler.Register)
		r.POST("/login", handler.Login)
		r.GET("/:userId", authMiddleware.ParseToken(), handler.GetUser)
		r.GET("/", authMiddleware.ParseToken(), handler.GetAllUser)
		r.PATCH("/:userId", authMiddleware.ParseToken(), handler.UpdateUser)
		r.DELETE("/:userId", authMiddleware.ParseToken(), handler.DeleteUser)
	}
}

// swagger:route POST /users/register users registerRequest
// Register a new user account.
// Responses:
// 	200: genericResponse
// 	400: genericError
func (h *UserHandler) Register(c *gin.Context) {
	var user *entities.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: err.Error()})
		return
	}

	err = h.UserUsecase.Register(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, swagger.GenericResponse{Status: "success"})
}

// swagger:route POST /users/login users loginRequest
// Log user into system.
// Responses:
// 	200: loginResponse
// 	400: genericError
func (h *UserHandler) Login(c *gin.Context) {
	var user *entities.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: err.Error()})
		return
	}

	userId, err := h.UserUsecase.Login(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: err.Error()})
		return
	}

	token, err := h.TokenUsecase.CreateToken(userId)
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

// swagger:operation GET /users/{userId} users getUser
// ---
// summary: Get the user profile by user id.
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
//     "$ref": "#/responses/getUserResponse"
//   "400":
//     "$ref": "#/responses/genericError"
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("userId")

	user, err := h.UserUsecase.GetUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: "invalid id"})
		return
	}

	c.JSON(http.StatusOK, &swagger.GetUserResponse{Status: "success", Data: user})
}

// swagger:route GET /users users getAllUser
// Get all the user profile.
// Responses:
// 	200: getAllUserResponse
// 	400: genericError
func (h *UserHandler) GetAllUser(c *gin.Context) {
	userList, err := h.UserUsecase.GetAllUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, swagger.GetAllUserResponse{Status: "success", Data: userList})
}

// swagger:operation PATCH /users/{userId} users updateUser
// ---
// summary: Update the user profile.
// parameters:
// - name: userId
//   in: path
//   description: user id
//   type: string
//   required: true
// - name: body
//   in: body
//   description: request body
//   schema:
//     "$ref": "#/definitions/UpdateUserRequest"
// security:
// - Bearer: []
// responses:
//   "200":
//     "$ref": "#/responses/genericResponse"
//   "400":
//     "$ref": "#/responses/genericError"
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user *entities.User
	id := c.Param("userId")

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: err.Error()})
		return
	}

	err = h.UserUsecase.UpdateUser(id, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, swagger.GenericResponse{Status: "success"})
}

// swagger:operation DELETE /users/{userId} users deleteUser
// ---
// summary: Delete a user.
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
//     "$ref": "#/responses/genericResponse"
//   "400":
//     "$ref": "#/responses/genericError"
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("userId")

	err := h.UserUsecase.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, swagger.GenericError{Status: "failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, swagger.GenericResponse{Status: "success"})
}

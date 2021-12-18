package http

import (
	"GoCleanArchitecture/delivery/http/middlewares"
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
		r.PATCH("/:userId", authMiddleware.ParseToken(), handler.PatchToUpdateUser)
		r.DELETE("/:userId", authMiddleware.ParseToken(), handler.DeleteUser)
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user entities.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	err = h.UserUsecase.Register(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": nil})
}

func (h *UserHandler) Login(c *gin.Context) {
	var user entities.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	userId, err := h.UserUsecase.Login(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	token, err := h.TokenUsecase.CreateToken(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	data := map[string]string{
		"user_id":       userId,
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("userId")

	user, err := h.UserUsecase.GetUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": "invalid id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
}

func (h *UserHandler) GetAllUser(c *gin.Context) {
	userList, err := h.UserUsecase.GetAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": userList})
}

func (h *UserHandler) PatchToUpdateUser(c *gin.Context) {
	var user entities.User
	id := c.Param("userId")

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	err = h.UserUsecase.UpdateUser(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": nil})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("userId")

	err := h.UserUsecase.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": nil})
}

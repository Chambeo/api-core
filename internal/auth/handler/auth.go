package handler

import (
	"chambeo-api-core/internal/auth/models"
	"chambeo-api-core/internal/users/service"
	"chambeo-api-core/pkg/customError"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type AuthHandlerInterface interface {
	GenerateToken(c *gin.Context)
	RefreshToken(c *gin.Context)
	ValidateToken(c *gin.Context)
}

type AuthService interface {
	GenerateToken(email string, userId string) string
	ParseToken(tokenString string) *jwt.Token
}

type AuthHandler struct {
	authService AuthService
	userService service.UserService
}

func NewAuthHandler(authService AuthService, userService service.UserService) AuthHandlerInterface {
	return AuthHandler{authService: authService, userService: userService}
}

func (a AuthHandler) GenerateToken(c *gin.Context) {
	var userDto models.UserLogin
	err := c.ShouldBindJSON(&userDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.Error{
			Code:    customError.InvalidBody,
			Message: "Invalid request body",
		})
		return
	}

	//user, err := a.userService.Get()

}

func (a AuthHandler) RefreshToken(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a AuthHandler) ValidateToken(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

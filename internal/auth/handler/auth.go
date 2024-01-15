package handler

import (
	"chambeo-api-core/internal/auth/models"
	"chambeo-api-core/internal/users/service"
	"chambeo-api-core/pkg/customError"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type AuthHandlerInterface interface {
	GenerateToken(c *gin.Context)
	RefreshToken(c *gin.Context)
	ValidateToken(c *gin.Context)
}

type AuthService interface {
	GenerateToken(email string, userId string) (*string, error)
	ParseToken(tokenString string) (*jwt.Token, error)
}

type AuthHandler struct {
	authService AuthService
	userService service.UserServiceInterface
}

func NewAuthHandler(authService AuthService, userService service.UserServiceInterface) AuthHandlerInterface {
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

	user, err := a.userService.GetByEmail(userDto.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: "Error trying to retrieve user from DB",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, customError.Error{
			Code:    customError.NotFound,
			Message: "User not found",
		})
		return
	}

	if !a.validPassword(userDto.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, customError.Error{
			Code:    customError.ApplicationError,
			Message: "Invalid credentials",
		})
		return
	}

	token, err := a.authService.GenerateToken(user.Email, strconv.Itoa(user.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: "Error trying to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{AccessToken: *token})
	return
}

func (a AuthHandler) RefreshToken(c *gin.Context) {
	var tokenToValidate models.TokenRequest
	err := c.ShouldBindJSON(&tokenToValidate)
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.Error{
			Code:    customError.InvalidBody,
			Message: "Invalid request body",
		})
		return
	}

	jwt, err := a.authService.ParseToken(tokenToValidate.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: "Error trying to parse token",
		})
		return
	}

	if !jwt.Valid {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: "Error trying to parse token",
		})
		return
	}

	claims, ok := jwt.Claims.(*models.CustomClaims)

	if !ok {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: "Error trying to parse token claims",
		})
		return
	}

	refreshedToken, err := a.authService.GenerateToken(claims.UserID, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: "Error trying to refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{AccessToken: *refreshedToken})
	return

}

func (a AuthHandler) ValidateToken(c *gin.Context) {
	var tokenToValidate models.TokenRequest
	err := c.ShouldBindJSON(&tokenToValidate)
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.Error{
			Code:    customError.InvalidBody,
			Message: "Invalid request body",
		})
		return
	}

	jwt, err := a.authService.ParseToken(tokenToValidate.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: "Error trying to parse token",
		})
		return
	}

	if !jwt.Valid {
		c.JSON(http.StatusUnauthorized, customError.Error{
			Code:    customError.ApplicationError,
			Message: "Token is not valid",
		})
		return
	}
	return
}

func (a AuthHandler) validPassword(requestPassword, retrievedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(retrievedPassword), []byte(requestPassword))
	if err != nil {
		return false
	}
	return true
}

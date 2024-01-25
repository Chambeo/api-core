package handler

import (
	"chambeo-api-core/internal/users/models"
	"chambeo-api-core/internal/users/service"
	"chambeo-api-core/pkg/customError"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandlerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetByEmail(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) UserHandlerInterface {
	return &UserHandler{userService}
}
func (u *UserHandler) Create(c *gin.Context) {
	var userDto models.UserRequest
	err := c.ShouldBindJSON(&userDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.Error{
			Code:    customError.InvalidBody,
			Message: "Invalid request body",
		})
		return
	}
	user, err := u.userService.Create(&userDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: "An error occurred when tyring to create user",
		})
		return
	}
	c.JSON(http.StatusCreated, user)
	return
}

func (u *UserHandler) Get(c *gin.Context) {
	userId := c.Param("id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, customError.Error{
			Code:    customError.MissingParameter,
			Message: "Missing or mismatch userId",
		})
		return
	}

	user, err := u.userService.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: fmt.Sprintf("An error occurred when trying to retrieve user with id %s", userId),
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

	c.JSON(http.StatusOK, user)
	return
}

func (u *UserHandler) Update(c *gin.Context) {
	var userDto models.UserRequest
	err := c.ShouldBindJSON(&userDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.Error{
			Code:    customError.InvalidBody,
			Message: "Invalid request body",
		})
		return
	}
	user, err := u.userService.Update(&userDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: "An error occurred when tyring to update user",
		})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

func (u *UserHandler) Delete(c *gin.Context) {
	userId := c.Param("id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, customError.Error{
			Code:    customError.MissingParameter,
			Message: "Missing or mismatch userId",
		})
		return
	}

	user, err := u.userService.Delete(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: fmt.Sprintf("An error occurred when trying to delete user with id %s", userId),
		})
		return
	}
	c.JSON(http.StatusNoContent, user)
	return
}

func (u *UserHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, customError.Error{
			Code:    customError.MissingParameter,
			Message: "Missing or mismatch email",
		})
		return
	}

	user, err := u.userService.GetByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: fmt.Sprintf("An error occurred when trying to retrieve user with email %s", email),
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

	c.JSON(http.StatusOK, user)
	return
}

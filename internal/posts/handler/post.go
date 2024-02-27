package handler

import (
	"chambeo-api-core/internal/posts/models"
	"chambeo-api-core/internal/posts/service"
	"chambeo-api-core/pkg/customError"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostHandlerInterface interface {
	Create(c *gin.Context)
	GetPosts(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type PostHandler struct {
	PostService service.PostServiceInterface
}

func NewPostHandler(postService service.PostServiceInterface) PostHandlerInterface {
	return &PostHandler{PostService: postService}
}

func (p *PostHandler) Create(c *gin.Context) {
	var postDto models.CreatePostRequest
	err := c.ShouldBindJSON(&postDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, customError.Error{
			Code:    customError.InvalidBody,
			Message: "Invalid request body",
		})
		return
	}

	post, err := p.PostService.Create(&postDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customError.Error{
			Code:    customError.ApplicationError,
			Message: "An error occurred when tyring to create post",
		})
		return
	}
	c.JSON(http.StatusCreated, post)
	return
}

func (p *PostHandler) GetPosts(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (p *PostHandler) Update(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (p *PostHandler) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

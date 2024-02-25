package handler

import (
	"chambeo-api-core/internal/posts/service"
	"github.com/gin-gonic/gin"
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

func (p PostHandler) Create(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (p PostHandler) GetPosts(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (p PostHandler) Update(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (p PostHandler) Delete(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewPostHandler(postService service.PostServiceInterface) PostHandlerInterface {
	return &PostHandler{PostService: postService}
}

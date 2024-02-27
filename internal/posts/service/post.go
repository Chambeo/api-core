package service

import "chambeo-api-core/internal/posts/models"

type PostServiceInterface interface {
	Create(postRequest *models.CreatePostRequest) (*models.CreatePostResult, error)
	Get()
	Update()
	Delete()
}

type PostService struct{}

func NewPostService() PostServiceInterface {
	return &PostService{}
}

func (p *PostService) Create(_ *models.CreatePostRequest) (*models.CreatePostResult, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostService) Get() {
	//TODO implement me
	panic("implement me")
}

func (p *PostService) Update() {
	//TODO implement me
	panic("implement me")
}

func (p *PostService) Delete() {
	//TODO implement me
	panic("implement me")
}

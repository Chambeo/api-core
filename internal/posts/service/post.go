package service

type PostServiceInterface interface {
}

type PostService struct{}

func NewPostService() PostServiceInterface {
	return &PostService{}
}

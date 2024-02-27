package handler

import (
	"bytes"
	"chambeo-api-core/internal/posts/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostHandler_Create(t *testing.T) {

	validPostDto := &models.CreatePostResult{
		ID:          "1",
		Title:       "Cuidado de perro",
		Amount:      "7000",
		PostType:    "OTROS",
		Description: "Se necesita cuidar a un perrito chico todo el dia que debe ser buscado del barrio de Colegiales el sabado",
		Status:      "PUBLISHED",
	}

	tests := []struct {
		name           string
		requestBody    string
		responseBody   string
		statusCode     int
		mockedBehavior func(postServiceMock *mock.Mock)
		asserts        func(t *testing.T, responseBody, receivedResponse string, statusCode, receivedStatusCode int)
	}{
		{
			name:         "a valid request body should be created successfully",
			requestBody:  `{"title":"Cuidado de perro","user_owner_id":"1","amount":"7000","post_type":"OTROS","description":"Se necesita cuidar a un perrito chico todo el dia que debe ser buscado del barrio de Colegiales el sabado"}`,
			responseBody: `{"id":"1","title":"Cuidado de perro","amount":"7000","post_type":"OTROS","description":"Se necesita cuidar a un perrito chico todo el dia que debe ser buscado del barrio de Colegiales el sabado","status":"PUBLISHED"}`,
			statusCode:   http.StatusCreated,
			mockedBehavior: func(postServiceMock *mock.Mock) {
				postServiceMock.On("Create", mock.Anything).Return(validPostDto, nil)
			},
			asserts: func(t *testing.T, responseBody, receivedResponse string, statusCode, receivedStatusCode int) {
				assert.Equal(t, responseBody, receivedResponse)
				assert.Equal(t, statusCode, receivedStatusCode)
			},
		},
		{
			name:         "a valid request body should return error from service",
			requestBody:  `{"title":"Cuidado de perro","user_owner_id":"1","amount":"7000","post_type":"OTROS","description":"Se necesita cuidar a un perrito chico todo el dia que debe ser buscado del barrio de Colegiales el sabado"}`,
			responseBody: `{"code":"ERROR","message":"An error occurred when tyring to create post"}`,
			statusCode:   http.StatusInternalServerError,
			mockedBehavior: func(postServiceMock *mock.Mock) {
				postServiceMock.On("Create", mock.Anything).Return(nil, errors.New("error"))
			},
			asserts: func(t *testing.T, responseBody, receivedResponse string, statusCode, receivedStatusCode int) {
				assert.Equal(t, responseBody, receivedResponse)
				assert.Equal(t, statusCode, receivedStatusCode)
			},
		},
		{
			name:         "an invalid request body should return error",
			requestBody:  `{`,
			responseBody: `{"code":"INVALID_BODY","message":"Invalid request body"}`,
			statusCode:   http.StatusBadRequest,
			mockedBehavior: func(postServiceMock *mock.Mock) {
				postServiceMock.On("Create", mock.Anything).Return(nil, errors.New("error"))
			},
			asserts: func(t *testing.T, responseBody, receivedResponse string, statusCode, receivedStatusCode int) {
				assert.Equal(t, responseBody, receivedResponse)
				assert.Equal(t, statusCode, receivedStatusCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			postService := PostServiceMock{}

			postHandler := NewPostHandler(&postService)

			tt.mockedBehavior(&postService.Mock)

			r := setupRouter(postHandler)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/posts/", bytes.NewReader([]byte(tt.requestBody)))

			r.ServeHTTP(w, req)

			tt.asserts(t, tt.responseBody, w.Body.String(), tt.statusCode, w.Code)

		})
	}

}

func setupRouter(postHandler PostHandlerInterface) *gin.Engine {
	r := gin.Default()
	r.POST("/posts/", postHandler.Create)
	return r
}

type PostServiceMock struct {
	mock.Mock
}

func (m *PostServiceMock) Create(postRequest *models.CreatePostRequest) (*models.CreatePostResult, error) {
	args := m.Called(postRequest)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CreatePostResult), nil
}

func (m *PostServiceMock) Get()    {}
func (m *PostServiceMock) Update() {}
func (m *PostServiceMock) Delete() {}

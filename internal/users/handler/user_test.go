package handler

import (
	"bytes"
	"chambeo-api-core/internal/users/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_Create(t *testing.T) {

	createdAt, _ := time.Parse(time.RFC3339, "2024-01-05T23:01:41.9180793-03:00")
	updatedAt, _ := time.Parse(time.RFC3339, "2024-01-05T23:01:41.9180793-03:00")

	tests := []struct {
		name                       string
		requestBody                string
		expectedBodyResponse       string
		expectedHttpStatusResponse int
		mockedBehavior             func(t *testing.T, mock *mock.Mock)
		asserts                    func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string)
	}{
		{
			name: "Valid request body should return 201 status code",
			requestBody: `{
					  "id": 1,
					  "first_name": "Meze",
					  "last_name": "Lawyer",
					  "email": "meze@email.com",
					  "password": "password"
					}
					`,
			expectedBodyResponse:       `{"id":1,"first_name":"Meze","last_name":"Lawyer","email":"meze@email.com","password":"password","created_at":"2024-01-05T23:01:41.9180793-03:00","updated_at":"2024-01-05T23:01:41.9180793-03:00"}`,
			expectedHttpStatusResponse: http.StatusCreated,
			mockedBehavior: func(t *testing.T, mockedService *mock.Mock) {
				mockedService.On("Create", mock.Anything).Return(&models.UserDto{
					Id:        1,
					FirstName: "Meze",
					LastName:  "Lawyer",
					Email:     "meze@email.com",
					Password:  "password",
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
					DeletedAt: nil,
				}, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "Test with invalid body should return 400",
			requestBody:                `}{,`,
			expectedBodyResponse:       "",
			mockedBehavior:             func(t *testing.T, mockedService *mock.Mock) {},
			expectedHttpStatusResponse: http.StatusBadRequest,
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
			},
		},
		{
			name: "Test with valid data should return 500 due service error",
			requestBody: `{
					  "id": 1,
					  "first_name": "Meze",
					  "last_name": "Lawyer",
					  "email": "meze@email.com",
					  "password": "password"
					}
					`,
			expectedBodyResponse:       `{"code":"ERROR","message":"An error occurred when tyring to create user"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
			mockedBehavior: func(t *testing.T, mockedService *mock.Mock) {
				mockedService.On("Create", mock.Anything).Return(nil, errors.New("error from service"))
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockedService := &MockUserService{}

			tt.mockedBehavior(t, &mockedService.Mock)

			userHandler := NewUserHandler(mockedService)

			router := setupMockedRouter(userHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewReader([]byte(tt.requestBody)))

			router.ServeHTTP(w, req)

			tt.asserts(t, w, tt.expectedHttpStatusResponse, tt.expectedBodyResponse)

		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2024-01-05T23:01:41.9180793-03:00")
	updatedAt, _ := time.Parse(time.RFC3339, "2024-01-05T23:01:41.9180793-03:00")

	tests := []struct {
		name                       string
		requestBody                string
		expectedBodyResponse       string
		expectedHttpStatusResponse int
		mockedBehavior             func(t *testing.T, mock *mock.Mock)
		asserts                    func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string)
	}{
		{
			name: "Valid request body should return 200 status code",
			requestBody: `{
					  "id": 1,
					  "first_name": "Meze",
					  "last_name": "Lawyer",
					  "email": "meze@email.com",
					  "password": "password"
					}
					`,
			expectedBodyResponse:       `{"id":1,"first_name":"Meze","last_name":"Lawyer","email":"meze@email.com","password":"password","created_at":"2024-01-05T23:01:41.9180793-03:00","updated_at":"2024-01-05T23:01:41.9180793-03:00"}`,
			expectedHttpStatusResponse: http.StatusOK,
			mockedBehavior: func(t *testing.T, mockedService *mock.Mock) {
				mockedService.On("Update", mock.Anything).Return(&models.UserDto{
					Id:        1,
					FirstName: "Meze",
					LastName:  "Lawyer",
					Email:     "meze@email.com",
					Password:  "password",
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
					DeletedAt: nil,
				}, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "Test with invalid body should return 400",
			requestBody:                `}{,`,
			expectedBodyResponse:       `{"code":"INVALID_BODY","message":"Invalid request body"}`,
			mockedBehavior:             func(t *testing.T, mockedService *mock.Mock) {},
			expectedHttpStatusResponse: http.StatusBadRequest,
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name: "Test with valid data should return 500 due service error",
			requestBody: `{
					  "id": 1,
					  "first_name": "Meze",
					  "last_name": "Lawyer",
					  "email": "meze@email.com",
					  "password": "password"
					}
					`,
			expectedBodyResponse:       `{"code":"ERROR","message":"An error occurred when tyring to update user"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
			mockedBehavior: func(t *testing.T, mockedService *mock.Mock) {
				mockedService.On("Update", mock.Anything).Return(nil, errors.New("error from service"))
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockedService := &MockUserService{}

			tt.mockedBehavior(t, &mockedService.Mock)

			userHandler := NewUserHandler(mockedService)

			router := setupMockedRouter(userHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/api/v1/users/", bytes.NewReader([]byte(tt.requestBody)))

			router.ServeHTTP(w, req)

			tt.asserts(t, w, tt.expectedHttpStatusResponse, tt.expectedBodyResponse)

		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2024-01-05T23:01:41.9180793-03:00")
	updatedAt, _ := time.Parse(time.RFC3339, "2024-01-05T23:01:41.9180793-03:00")

	tests := []struct {
		name                       string
		id                         string
		expectedBodyResponse       string
		expectedHttpStatusResponse int
		mockedBehavior             func(t *testing.T, mock *mock.Mock)
		asserts                    func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string)
	}{
		{
			name:                       "Valid request should return 200 status code",
			id:                         "1",
			expectedBodyResponse:       `{"id":1,"first_name":"Meze","last_name":"Lawyer","email":"meze@email.com","password":"password","created_at":"2024-01-05T23:01:41.9180793-03:00","updated_at":"2024-01-05T23:01:41.9180793-03:00"}`,
			expectedHttpStatusResponse: http.StatusOK,
			mockedBehavior: func(t *testing.T, mockedService *mock.Mock) {
				mockedService.On("Get", mock.Anything).Return(&models.UserDto{
					Id:        1,
					FirstName: "Meze",
					LastName:  "Lawyer",
					Email:     "meze@email.com",
					Password:  "password",
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
					DeletedAt: nil,
				}, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                 "Test with valid data should return 404",
			id:                   "1",
			expectedBodyResponse: `{"code":"NOT_FOUND","message":"User not found"}`,
			mockedBehavior: func(t *testing.T, mockedService *mock.Mock) {
				mockedService.On("Get", mock.Anything).Return(nil, nil)
			},
			expectedHttpStatusResponse: http.StatusNotFound,
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "Test with valid data should return 500 due service error",
			id:                         "1",
			expectedBodyResponse:       `{"code":"ERROR","message":"An error occurred when trying to retrieve user with id 1"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
			mockedBehavior: func(t *testing.T, mockedService *mock.Mock) {
				mockedService.On("Get", mock.Anything).Return(nil, errors.New("error from service"))
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockedService := &MockUserService{}

			tt.mockedBehavior(t, &mockedService.Mock)

			userHandler := NewUserHandler(mockedService)

			router := setupMockedRouter(userHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", tt.id), nil)

			router.ServeHTTP(w, req)

			tt.asserts(t, w, tt.expectedHttpStatusResponse, tt.expectedBodyResponse)

		})
	}
}

func TestUserHandler_Delete(t *testing.T) {

	tests := []struct {
		name                       string
		id                         string
		expectedBodyResponse       string
		expectedHttpStatusResponse int
		mockedBehavior             func(t *testing.T, mock *mock.Mock)
		asserts                    func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string)
	}{
		{
			name:                       "Valid request should return 204 status code",
			id:                         "1",
			expectedBodyResponse:       "",
			expectedHttpStatusResponse: http.StatusNoContent,
			mockedBehavior: func(t *testing.T, mockedService *mock.Mock) {
				mockedService.On("Delete", mock.Anything).Return(&models.UserDto{}, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "Test with valid data should return 500 due service error",
			id:                         "1",
			expectedBodyResponse:       `{"code":"ERROR","message":"An error occurred when trying to delete user with id 1"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
			mockedBehavior: func(t *testing.T, mockedService *mock.Mock) {
				mockedService.On("Delete", mock.Anything).Return(nil, errors.New("error from service"))
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockedService := &MockUserService{}

			tt.mockedBehavior(t, &mockedService.Mock)

			userHandler := NewUserHandler(mockedService)

			router := setupMockedRouter(userHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%s", tt.id), nil)

			router.ServeHTTP(w, req)

			tt.asserts(t, w, tt.expectedHttpStatusResponse, tt.expectedBodyResponse)

		})
	}
}

func setupMockedRouter(userHandler UserHandlerInterface) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/", userHandler.Create)
			users.GET("/:id", userHandler.Get)
			users.PUT("/", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}

	}

	return r
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Get(id string) (*models.UserDto, error) {
	args := m.Called(id)
	if args.Get(1) != nil || args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserDto), args.Error(1)
}

func (m *MockUserService) Update(user *models.UserDto) (*models.UserDto, error) {
	args := m.Called(user)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserDto), args.Error(1)
}

func (m *MockUserService) Delete(id string) (*models.UserDto, error) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserDto), args.Error(1)
}

func (m *MockUserService) Create(user *models.UserDto) (*models.UserDto, error) {
	args := m.Called(user)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.UserDto), args.Error(1)
}

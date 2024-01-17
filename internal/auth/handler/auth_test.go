package handler

import (
	"bytes"
	"chambeo-api-core/internal/users/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthHandler_GenerateToken(t *testing.T) {

	mockedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsImVtYWlsIjoibWV6ZUBnbWFpb" +
		"C5jb20iLCJpc3MiOiJjaGFtYmVvLWNvIiwic3ViIjoiY2hhbWJlby1iZSIsImF1ZCI6WyJjaGFtYmVvLWZlIl0sImV4" +
		"cCI6MTcwNTI3NjMyMiwibmJmIjoxNzA1MTg5OTIyLCJpYXQiOjE3MDUxODk5MjIsImp0aSI6IjEifQ.p2jndX8Bn8q3" +
		"mrJp4vv9nsGugZOZRcukrOBuMSIO4SA"

	createdAt := time.Now()
	updatedAt := createdAt

	tests := []struct {
		name                       string
		requestBody                string
		expectedBodyResponse       string
		expectedHttpStatusResponse int
		mockedBehavior             func(t *testing.T, userMock, authMock *mock.Mock)
		asserts                    func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string)
	}{
		{
			name:                       "valid credentials should return token",
			requestBody:                `{"email":"meze@gmail.com", "password":"password"}`,
			expectedBodyResponse:       "",
			expectedHttpStatusResponse: http.StatusOK,
			mockedBehavior: func(t *testing.T, userMock, authMock *mock.Mock) {
				userMock.On("GetByEmail", mock.Anything).Return(&models.UserDto{
					Id:        1,
					FirstName: "Meze",
					LastName:  "Lawyer",
					Email:     "meze@gmail.com",
					Password:  "$2a$16$m.fWPulWk20mcpq5lZnkMeB7sOu2w10o/3EGwjLURZ3A7AcI9O4lC",
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
					DeletedAt: nil,
				}, nil)
				authMock.On("GenerateToken", "meze@gmail.com", "1").Return(&mockedToken, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, fmt.Sprintf(`{"access_token":"%s"}`, mockedToken), response.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockedUserService := &MockUserService{}
			mockedAuthService := &MockAuthService{}

			tt.mockedBehavior(t, &mockedUserService.Mock, &mockedAuthService.Mock)

			authHandler := NewAuthHandler(mockedAuthService, mockedUserService)

			router := setupMockedRouter(authHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/auth/token", bytes.NewReader([]byte(tt.requestBody)))

			router.ServeHTTP(w, req)

			tt.asserts(t, w, tt.expectedHttpStatusResponse, tt.expectedBodyResponse)

		})
	}
}

func setupMockedRouter(authHandler AuthHandlerInterface) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/auth")
		{
			users.POST("/token", authHandler.GenerateToken)
			users.GET("/token/validate", authHandler.ValidateToken)
			users.POST("/token/refresh", authHandler.RefreshToken)
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

func (m *MockUserService) GetByEmail(email string) (*models.UserDto, error) {
	args := m.Called(email)
	if args.Get(1) != nil || args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserDto), args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GenerateToken(email, userId string) (*string, error) {
	args := m.Called(email, userId)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *MockAuthService) ParseToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

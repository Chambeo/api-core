package handler

import (
	"bytes"
	authClaims "chambeo-api-core/internal/auth/models"
	"chambeo-api-core/internal/users/models"
	"errors"
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
		{
			name:                       "invalid credentials should return error",
			requestBody:                `{"email":"meze@gmail.com"}`,
			expectedBodyResponse:       `{"code":"INVALID_BODY","message":"Invalid request body"}`,
			expectedHttpStatusResponse: http.StatusBadRequest,
			mockedBehavior:             func(t *testing.T, userMock, authMock *mock.Mock) {},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "valid credentials should return error getting user",
			requestBody:                `{"email":"meze@gmail.com", "password":"password"}`,
			expectedBodyResponse:       `{"code":"ERROR","message":"Error trying to retrieve user from DB"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
			mockedBehavior: func(t *testing.T, userMock, authMock *mock.Mock) {
				userMock.On("GetByEmail", mock.Anything).Return(nil, errors.New("error from repository"))
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "valid credentials should return user not found",
			requestBody:                `{"email":"meze@gmail.com", "password":"password"}`,
			expectedBodyResponse:       `{"code":"NOT_FOUND","message":"User not found"}`,
			expectedHttpStatusResponse: http.StatusNotFound,
			mockedBehavior: func(t *testing.T, userMock, authMock *mock.Mock) {
				userMock.On("GetByEmail", mock.Anything).Return(nil, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "invalid credentials should return unauthorized",
			requestBody:                `{"email":"meze@gmail.com", "password":"invalidPassword"}`,
			expectedBodyResponse:       `{"code":"ERROR","message":"Invalid credentials"}`,
			expectedHttpStatusResponse: http.StatusUnauthorized,
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
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "valid credentials should return error when trying to generate token",
			requestBody:                `{"email":"meze@gmail.com", "password":"password"}`,
			expectedBodyResponse:       `{"code":"ERROR","message":"Error trying to generate token"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
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
				authMock.On("GenerateToken", "meze@gmail.com", "1").Return(nil, errors.New("error generating token"))
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
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

func TestAuthHandler_RefreshToken(t *testing.T) {

	validTokenResponse := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsImVtYWlsIjoibWV6ZUBnbWFpbC5jb20iLCJpc3MiOiJjaG" +
		"FtYmVvLWNvIiwic3ViIjoiY2hhbWJlby1iZSIsImF1ZCI6WyJjaGFtYmVvLWZlIl0sImV4cCI6MTcwNTI3NjMyMiwibmJmIjoxNzA1MTg5OTI" +
		"yLCJpYXQiOjE3MDUxODk5MjIsImp0aSI6IjEifQ.p2jndX8Bn8q3mrJp4vv9nsGugZOZRcukrOBuMSIO4SA"

	validJwtParse := &jwt.Token{
		Raw:    validTokenResponse,
		Method: nil,
		Header: nil,
		Claims: &authClaims.CustomClaims{
			UserID: "1",
			Email:  "meze@gmail.com",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "",
				Subject:   "",
				Audience:  nil,
				ExpiresAt: nil,
				NotBefore: nil,
				IssuedAt:  nil,
				ID:        "",
			},
		},
		Signature: nil,
		Valid:     true,
	}

	validJwtParseWithInvalidClaims := &jwt.Token{
		Raw:       validTokenResponse,
		Method:    nil,
		Header:    nil,
		Claims:    jwt.RegisteredClaims{},
		Signature: nil,
		Valid:     true,
	}

	invalidJwtParseWithInvalidClaims := &jwt.Token{
		Raw:       validTokenResponse,
		Method:    nil,
		Header:    nil,
		Claims:    jwt.RegisteredClaims{},
		Signature: nil,
		Valid:     false,
	}

	tests := []struct {
		name                       string
		requestBody                string
		expectedBodyResponse       string
		expectedHttpStatusResponse int
		mockedBehavior             func(t *testing.T, authMock *mock.Mock)
		asserts                    func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string)
	}{
		{
			name:                       "valid token should be refreshed",
			requestBody:                fmt.Sprintf(`{"access_token":"%s"}`, validTokenResponse),
			expectedBodyResponse:       fmt.Sprintf(`{"access_token":"%s"}`, validTokenResponse),
			expectedHttpStatusResponse: http.StatusOK,
			mockedBehavior: func(t *testing.T, authMock *mock.Mock) {
				authMock.On("ParseToken", mock.Anything).Return(validJwtParse, nil)
				authMock.On("GenerateToken", mock.Anything, mock.Anything).Return(&validTokenResponse, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "valid token return error when refreshing token",
			requestBody:                fmt.Sprintf(`{"access_token":"%s"}`, validTokenResponse),
			expectedBodyResponse:       `{"code":"ERROR","message":"Error trying to refresh token"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
			mockedBehavior: func(t *testing.T, authMock *mock.Mock) {
				authMock.On("ParseToken", mock.Anything).Return(validJwtParse, nil)
				authMock.On("GenerateToken", mock.Anything, mock.Anything).Return(nil, errors.New("error refreshing token"))
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "valid token return error due invalid claims",
			requestBody:                fmt.Sprintf(`{"access_token":"%s"}`, validTokenResponse),
			expectedBodyResponse:       `{"code":"ERROR","message":"Error trying to parse token claims"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
			mockedBehavior: func(t *testing.T, authMock *mock.Mock) {
				authMock.On("ParseToken", mock.Anything).Return(validJwtParseWithInvalidClaims, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "valid token return error due invalid jwt refresh",
			requestBody:                fmt.Sprintf(`{"access_token":"%s"}`, validTokenResponse),
			expectedBodyResponse:       `{"code":"ERROR","message":"Error refreshing token"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
			mockedBehavior: func(t *testing.T, authMock *mock.Mock) {
				authMock.On("ParseToken", mock.Anything).Return(invalidJwtParseWithInvalidClaims, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "invalid token return error due error parsing",
			requestBody:                fmt.Sprintf(`{"access_token":"%s"}`, validTokenResponse),
			expectedBodyResponse:       `{"code":"ERROR","message":"Error trying to parse token"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
			mockedBehavior: func(t *testing.T, authMock *mock.Mock) {
				authMock.On("ParseToken", mock.Anything).Return(nil, errors.New("error parsing token"))
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
		{
			name:                       "invalid request token return bad request",
			requestBody:                `{"name":"meze"}`,
			expectedBodyResponse:       `{"code":"INVALID_BODY","message":"Invalid request body"}`,
			expectedHttpStatusResponse: http.StatusBadRequest,
			mockedBehavior:             func(t *testing.T, authMock *mock.Mock) {},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
				assert.Equal(t, expectedBody, response.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedAuthService := &MockAuthService{}
			tt.mockedBehavior(t, &mockedAuthService.Mock)

			authHandler := NewAuthHandler(mockedAuthService, nil)

			router := setupMockedRouter(authHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/auth/token/refresh", bytes.NewReader([]byte(tt.requestBody)))

			router.ServeHTTP(w, req)

			tt.asserts(t, w, tt.expectedHttpStatusResponse, tt.expectedBodyResponse)
		})
	}
}

func TestAuthHandler_ValidateToken(t *testing.T) {

	validJwtParse := &jwt.Token{
		Raw:    "validToken",
		Method: nil,
		Header: nil,
		Claims: &authClaims.CustomClaims{
			UserID: "1",
			Email:  "meze@gmail.com",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "",
				Subject:   "",
				Audience:  nil,
				ExpiresAt: nil,
				NotBefore: nil,
				IssuedAt:  nil,
				ID:        "",
			},
		},
		Signature: nil,
		Valid:     true,
	}

	invalidJwtParse := &jwt.Token{
		Raw:    "validToken",
		Method: nil,
		Header: nil,
		Claims: &authClaims.CustomClaims{
			UserID: "1",
			Email:  "meze@gmail.com",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "",
				Subject:   "",
				Audience:  nil,
				ExpiresAt: nil,
				NotBefore: nil,
				IssuedAt:  nil,
				ID:        "",
			},
		},
		Signature: nil,
		Valid:     false,
	}

	tests := []struct {
		name                       string
		requestBody                string
		expectedBodyResponse       string
		expectedHttpStatusResponse int
		mockedBehavior             func(t *testing.T, authMock *mock.Mock)
		asserts                    func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string)
	}{
		{
			name:                       "invalid body should return bad request",
			requestBody:                `{"email":"meze@mail.com"}`,
			expectedBodyResponse:       `{"code":"INVALID_BODY","message":"Invalid request body"}`,
			expectedHttpStatusResponse: http.StatusBadRequest,
			mockedBehavior:             func(t *testing.T, authMock *mock.Mock) {},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedBody, response.Body.String())
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
			},
		},
		{
			name:                       "valid body should return error when parsing",
			requestBody:                `{"access_token":"exampleMockedToken"}`,
			expectedBodyResponse:       `{"code":"ERROR","message":"Error trying to parse token"}`,
			expectedHttpStatusResponse: http.StatusInternalServerError,
			mockedBehavior: func(t *testing.T, authMock *mock.Mock) {
				authMock.On("ParseToken", mock.Anything).Return(nil, errors.New("error parsing token"))
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedBody, response.Body.String())
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
			},
		},
		{
			name:                       "refreshed token is not valid",
			requestBody:                `{"access_token":"exampleMockedToken"}`,
			expectedBodyResponse:       `{"code":"ERROR","message":"Token is not valid"}`,
			expectedHttpStatusResponse: http.StatusUnauthorized,
			mockedBehavior: func(t *testing.T, authMock *mock.Mock) {
				authMock.On("ParseToken", mock.Anything).Return(invalidJwtParse, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedBody, response.Body.String())
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
			},
		},
		{
			name:                       "refreshed token is valid",
			requestBody:                `{"access_token":"exampleMockedToken"}`,
			expectedBodyResponse:       `{"access_token":"validToken"}`,
			expectedHttpStatusResponse: http.StatusOK,
			mockedBehavior: func(t *testing.T, authMock *mock.Mock) {
				authMock.On("ParseToken", mock.Anything).Return(validJwtParse, nil)
			},
			asserts: func(t *testing.T, response *httptest.ResponseRecorder, expectedHttpStatusResponse int, expectedBody string) {
				assert.Equal(t, expectedHttpStatusResponse, response.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedAuthService := &MockAuthService{}
			tt.mockedBehavior(t, &mockedAuthService.Mock)

			authHandler := NewAuthHandler(mockedAuthService, nil)

			router := setupMockedRouter(authHandler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/auth/token/validate", bytes.NewReader([]byte(tt.requestBody)))

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

package service

import (
	"chambeo-api-core/internal/users/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
	"time"
)

var (
	validUserRequest = &models.UserRequest{
		Id:        1,
		FirstName: "Meze",
		LastName:  "Lawyer",
		Email:     "meze@gmail.com",
		Password:  "password",
	}

	validUserResponse = &models.UserRequest{
		Id:        1,
		FirstName: "Meze",
		LastName:  "Lawyer",
		Email:     "meze@gmail.com",
		Password:  "$2a$16$xYdZVlrHaga5BC6Hw80a/u9QRl5BMJXew8bBKlzA3mpVPgVwk.WMi",
		DeletedAt: &time.Time{},
	}

	validUserModel = &models.User{
		Model: gorm.Model{
			ID:        1,
			DeletedAt: gorm.DeletedAt{},
		},
		FirstName: "Meze",
		LastName:  "Lawyer",
		Email:     "meze@gmail.com",
		Password:  "$2a$16$xYdZVlrHaga5BC6Hw80a/u9QRl5BMJXew8bBKlzA3mpVPgVwk.WMi",
	}
)

func TestUserService_Create(t *testing.T) {

	tests := []struct {
		name           string
		mockedBehavior func(t *testing.T, mockedRepository *mock.Mock)
		asserts        func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error)
		request        *models.UserRequest
		response       *models.UserRequest
		error          error
	}{
		{
			name: "With valid data should pass successfully",
			mockedBehavior: func(t *testing.T, mockedRepository *mock.Mock) {
				mockedRepository.On("Create", mock.Anything).Return(validUserModel, nil)
			},
			asserts: func(t *testing.T, response *models.UserRequest, error error, expectedError error) {
				assert.NotNil(t, response)
				assert.Nil(t, error)
				assert.Equal(t, validUserResponse, response)
			},
			request:  validUserRequest,
			response: validUserResponse,
			error:    nil,
		},
		{
			name: "Test with valid data should fail due repository error",
			mockedBehavior: func(t *testing.T, mockedRepository *mock.Mock) {
				mockedRepository.On("Create", mock.Anything).Return(nil, errors.New("error from repo"))
			},
			asserts: func(t *testing.T, response *models.UserRequest, error error, expectedError error) {
				assert.Nil(t, response)
				assert.NotNil(t, error)
				assert.Equal(t, error.Error(), expectedError.Error())
			},
			request:  validUserRequest,
			response: nil,
			error:    errors.New("error from repo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			userRepository := &MockUserRepository{}

			tt.mockedBehavior(t, &userRepository.Mock)

			userService := NewUser(userRepository)

			result, err := userService.Create(tt.request)

			tt.asserts(t, result, err, tt.error)

		})
	}

}

func TestUserService_Get(t *testing.T) {

	tests := []struct {
		name           string
		id             string
		mockedBehavior func(t *testing.T, mockedRepository *mock.Mock)
		asserts        func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error)
		response       *models.UserRequest
		error          error
	}{
		{
			name: "Get by id should return results",
			id:   "1",
			mockedBehavior: func(t *testing.T, mockedRepository *mock.Mock) {
				mockedRepository.On("Get", mock.Anything).Return(validUserModel, nil)
			},
			asserts: func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error) {
				assert.NotNil(t, response)
				assert.Nil(t, errorResult)
			},
			response: validUserResponse,
			error:    nil,
		},
		{
			name: "Get by id should return error",
			id:   "1",
			mockedBehavior: func(t *testing.T, mockedRepository *mock.Mock) {
				mockedRepository.On("Get", mock.Anything).Return(nil, errors.New("error"))
			},
			asserts: func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error) {
				assert.Nil(t, response)
				assert.NotNil(t, errorResult)
				assert.Equal(t, errorResult.Error(), expectedError.Error())
			},
			response: nil,
			error:    errors.New("ocurrio un error al intentar recuperar el usuario"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepository := &MockUserRepository{}

			tt.mockedBehavior(t, &userRepository.Mock)

			userService := NewUser(userRepository)

			result, err := userService.Get(tt.id)

			tt.asserts(t, result, err, tt.error)
		})
	}
}

func TestUserService_GetByEmail(t *testing.T) {

	tests := []struct {
		name           string
		email          string
		mockedBehavior func(t *testing.T, mockedRepository *mock.Mock)
		asserts        func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error)
		response       *models.UserRequest
		error          error
	}{
		{
			name:  "Get by email should return results",
			email: "meze@email.com",
			mockedBehavior: func(t *testing.T, mockedRepository *mock.Mock) {
				mockedRepository.On("GetByEmail", mock.Anything).Return(validUserModel, nil)
			},
			asserts: func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error) {
				assert.NotNil(t, response)
				assert.Nil(t, errorResult)
			},
			response: validUserResponse,
			error:    nil,
		},
		{
			name:  "Get by email should return error",
			email: "meze@email.com",
			mockedBehavior: func(t *testing.T, mockedRepository *mock.Mock) {
				mockedRepository.On("GetByEmail", mock.Anything).Return(nil, errors.New("error"))
			},
			asserts: func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error) {
				assert.Nil(t, response)
				assert.NotNil(t, errorResult)
				assert.Equal(t, errorResult.Error(), expectedError.Error())
			},
			response: nil,
			error:    errors.New("ocurrio un error al intentar recuperar el usuario"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepository := &MockUserRepository{}

			tt.mockedBehavior(t, &userRepository.Mock)

			userService := NewUser(userRepository)

			result, err := userService.GetByEmail(tt.email)

			tt.asserts(t, result, err, tt.error)
		})
	}

}

func TestUserService_Update(t *testing.T) {

	tests := []struct {
		name           string
		mockedBehavior func(t *testing.T, mockedRepository *mock.Mock)
		asserts        func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error)
		request        *models.UserRequest
		response       *models.UserRequest
		error          error
	}{
		{
			name: "With valid data should pass successfully",
			mockedBehavior: func(t *testing.T, mockedRepository *mock.Mock) {
				mockedRepository.On("Update", mock.Anything).Return(validUserModel, nil)
			},
			asserts: func(t *testing.T, response *models.UserRequest, error error, expectedError error) {
				assert.NotNil(t, response)
				assert.Nil(t, error)
				assert.Equal(t, validUserResponse, response)
			},
			request:  validUserRequest,
			response: validUserResponse,
			error:    nil,
		},
		{
			name: "Test with valid data should fail due repository error",
			mockedBehavior: func(t *testing.T, mockedRepository *mock.Mock) {
				mockedRepository.On("Update", mock.Anything).Return(nil, errors.New("error from repo"))
			},
			asserts: func(t *testing.T, response *models.UserRequest, error error, expectedError error) {
				assert.Nil(t, response)
				assert.NotNil(t, error)
				assert.Equal(t, error.Error(), expectedError.Error())
			},
			request:  validUserRequest,
			response: nil,
			error:    errors.New("ocurrio un error al intentar actualizar el usuario"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			userRepository := &MockUserRepository{}

			tt.mockedBehavior(t, &userRepository.Mock)

			userService := NewUser(userRepository)

			result, err := userService.Update(tt.request)

			tt.asserts(t, result, err, tt.error)

		})
	}

}

func TestUserService_Delete(t *testing.T) {

	tests := []struct {
		name           string
		id             string
		mockedBehavior func(t *testing.T, mockedRepository *mock.Mock)
		asserts        func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error)
		response       *models.UserRequest
		error          error
	}{
		{
			name: "Delete by id should return results",
			id:   "1",
			mockedBehavior: func(t *testing.T, mockedRepository *mock.Mock) {
				mockedRepository.On("Delete", mock.Anything).Return(validUserModel, nil)
			},
			asserts: func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error) {
				assert.NotNil(t, response)
				assert.Nil(t, errorResult)
			},
			response: validUserResponse,
			error:    nil,
		},
		{
			name: "Delete by id should return error",
			id:   "1",
			mockedBehavior: func(t *testing.T, mockedRepository *mock.Mock) {
				mockedRepository.On("Delete", mock.Anything).Return(nil, errors.New("error"))
			},
			asserts: func(t *testing.T, response *models.UserRequest, errorResult error, expectedError error) {
				assert.Nil(t, response)
				assert.NotNil(t, errorResult)
				assert.Equal(t, errorResult.Error(), expectedError.Error())
			},
			response: nil,
			error:    errors.New("ocurrio un error al intentar eliminar el usuario"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepository := &MockUserRepository{}

			tt.mockedBehavior(t, &userRepository.Mock)

			userService := NewUser(userRepository)

			result, err := userService.Delete(tt.id)

			tt.asserts(t, result, err, tt.error)
		})
	}

}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Get(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

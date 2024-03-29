package service

import (
	"chambeo-api-core/internal/users/models"
	"chambeo-api-core/internal/users/repository"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

type UserServiceInterface interface {
	Create(user *models.UserRequest) (*models.UserRequest, error)
	Get(id string) (*models.UserRequest, error)
	GetByEmail(id string) (*models.UserRequest, error)
	Update(user *models.UserRequest) (*models.UserRequest, error)
	Delete(id string) (*models.UserRequest, error)
}

type UserService struct {
	userRepository repository.UserRepositoryInterface
}

func NewUser(userRepository repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{userRepository: userRepository}
}

func (u *UserService) Create(user *models.UserRequest) (*models.UserRequest, error) {

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 16)

	if err != nil {
		log.Println("error when trying to encrypt password")
		return nil, errors.New("error al generar la contrasena para la cuenta")
	}
	user.Password = string(encryptedPassword)

	create, err := u.userRepository.Create(mapUserDtoToUserDb(*user))

	if err != nil {
		return nil, err
	}
	return mapUserDbToDto(*create), nil
}

func (u *UserService) Get(id string) (*models.UserRequest, error) {
	user, err := u.userRepository.Get(id)
	if err != nil {
		log.Println(fmt.Sprintf("error occurred trying to retrieve user with id %s", id))
		return nil, errors.New("ocurrio un error al intentar recuperar el usuario")
	}
	return mapUserDbToDto(*user), nil
}

func (u *UserService) Update(user *models.UserRequest) (*models.UserRequest, error) {
	updatedUser, err := u.userRepository.Update(mapUserDtoToUserDb(*user))
	if err != nil {
		log.Println(fmt.Sprintf("An error occurred trying to update user with id %v", user.Id))
		return nil, errors.New("ocurrio un error al intentar actualizar el usuario")
	}
	return mapUserDbToDto(*updatedUser), nil
}

func (u *UserService) Delete(id string) (*models.UserRequest, error) {
	user, err := u.userRepository.Delete(id)
	if err != nil {
		log.Println(fmt.Sprintf("error occurred trying to delete user with id %s", id))
		return nil, errors.New("ocurrio un error al intentar eliminar el usuario")
	}
	return mapUserDbToDto(*user), nil
}

func (u *UserService) GetByEmail(email string) (*models.UserRequest, error) {
	user, err := u.userRepository.GetByEmail(email)
	if err != nil {
		log.Println(fmt.Sprintf("error occurred trying to retrieve user with email %s", email))
		return nil, errors.New("ocurrio un error al intentar recuperar el usuario")
	}
	return mapUserDbToDto(*user), nil
}

func mapUserDtoToUserDb(user models.UserRequest) *models.User {
	return &models.User{
		Model: gorm.Model{
			ID:        uint(user.Id),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}
}

func mapUserDbToDto(user models.User) *models.UserRequest {
	return &models.UserRequest{
		Id:        int(user.Model.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: &user.DeletedAt.Time,
	}
}

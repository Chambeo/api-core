package repository

import (
	"chambeo-api-core/internal/users/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type UserRepositoryInterface interface {
	Create(user *models.User) (*models.User, error)
	Get(id string) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id string) (*models.User, error)
}

type UserRepository struct {
	DB gorm.DB
}

func NewUser(db gorm.DB) UserRepositoryInterface {
	return &UserRepository{DB: db}
}

func (u *UserRepository) Create(user *models.User) (*models.User, error) {
	result := u.DB.Create(&user) // pass pointer of data to Create

	//user.ID             // returns inserted data's primary key
	//result.Error        // returns error
	//result.RowsAffected // returns inserted records count

	if result.Error != nil {
		log.Println("Error on insert user: ", result.Error.Error()) // TODO
		return nil, errors.New("error al insertar el usuario en DB")
	}

	return user, nil
}

func (u *UserRepository) Get(id string) (*models.User, error) {
	var user *models.User
	if tx := u.DB.First(&user, id); tx.Error != nil {
		log.Println(fmt.Sprintf("error retrieving user with id %s %s", id, tx.Error.Error()))
		return nil, errors.New("error al recuperar el usuario en DB")
	}
	return user, nil
}

func (u *UserRepository) Update(user *models.User) (*models.User, error) {
	if tx := u.DB.Save(user); tx.Error != nil {
		log.Println(fmt.Sprintf("Error trying to update user with id %s", user.Id))
		return nil, errors.New("error al actualizar el usuario en DB")
	}
	return user, nil
}

func (u *UserRepository) Delete(id string) (*models.User, error) {
	if tx := u.DB.Where("id = ?", id).Delete(&models.User{}); tx.Error != nil {
		log.Println(fmt.Sprintf("Error trying to delete user with id %s", id))
		return nil, errors.New("error al intentar eliminar el usuario")
	}
	return &models.User{}, nil // TODO
}

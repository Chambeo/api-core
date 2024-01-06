package repository

import (
	"chambeo-api-core/internal/users/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"regexp"
	"testing"
	"time"
)

func TestUserRepository_Create(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now()

	tests := []struct {
		name           string
		userRequest    *models.User
		mockedBehavior func(t *testing.T, mock sqlmock.Sqlmock, user *models.User)
		asserts        func(t *testing.T, usr *models.User, error error)
	}{
		{
			name: "create user should be successfull",
			userRequest: &models.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				FirstName: "Meze",
				LastName:  "Lawyer",
				Email:     "meze@gmail.com",
				Password:  "password",
			},
			mockedBehavior: func(t *testing.T, mock sqlmock.Sqlmock, validUser *models.User) {

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`first_name`,`last_name`,`email`,`password`,`id`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(validUser.CreatedAt, validUser.UpdatedAt, nil, validUser.FirstName, validUser.LastName, validUser.Email, validUser.Password, validUser.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

			},
			asserts: func(t *testing.T, usr *models.User, error error) {
				assert.NoError(t, error)
			},
		},
		{
			name: "create user should return err",
			userRequest: &models.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				FirstName: "Meze",
				LastName:  "Lawyer",
				Email:     "meze@gmail.com",
				Password:  "password",
			},
			mockedBehavior: func(t *testing.T, mock sqlmock.Sqlmock, validUser *models.User) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`first_name`,`last_name`,`email`,`password`,`id`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(validUser.CreatedAt, validUser.UpdatedAt, nil, validUser.FirstName, validUser.LastName, validUser.Email, validUser.Password, validUser.ID).
					WillReturnError(errors.New("error from db"))
				mock.ExpectCommit()
			},
			asserts: func(t *testing.T, usr *models.User, error error) {
				assert.NotNil(t, error)
				assert.Equal(t, "error al insertar el usuario en DB", error.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDb, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})

			if err != nil {
				t.Error(err.Error())
			}
			gormDb.Debug()

			tt.mockedBehavior(t, mock, tt.userRequest)

			repository := NewUser(*gormDb)

			result, err := repository.Create(tt.userRequest)

			tt.asserts(t, result, err)

		})
	}
}

func TestUserRepository_Get(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockedBehavior func(t *testing.T, mock sqlmock.Sqlmock, id string)
		asserts        func(t *testing.T, user *models.User, err error)
	}{
		{
			name: "Test with valid id should return user result",
			id:   "1",
			mockedBehavior: func(t *testing.T, mock sqlmock.Sqlmock, id string) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name"}).
					AddRow(1, "Martin", "Lawyer")

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(id).
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
			asserts: func(t *testing.T, user *models.User, err error) {
				assert.NotNil(t, user)
				assert.Nil(t, err)
			},
		},
		{
			name: "Test with a valid id should return error from db",
			id:   "1",
			mockedBehavior: func(t *testing.T, mock sqlmock.Sqlmock, id string) {

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(id).
					WillReturnError(errors.New("error al recuperar el usuario en DB"))
				mock.ExpectCommit()
			},
			asserts: func(t *testing.T, user *models.User, err error) {
				assert.Nil(t, user)
				assert.NotNil(t, err)
			},
		},
		{
			name: "Test with a valid id should return no data from db",
			id:   "1",
			mockedBehavior: func(t *testing.T, mock sqlmock.Sqlmock, id string) {
				emptyRows := &sqlmock.Rows{}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).WithArgs(id).
					WillReturnRows(emptyRows)
				mock.ExpectCommit()
			},
			asserts: func(t *testing.T, user *models.User, err error) {
				assert.Nil(t, user)
				assert.NotNil(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDb, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})

			if err != nil {
				t.Error(err.Error())
			}
			gormDb.Debug()

			tt.mockedBehavior(t, mock, tt.id)

			repository := NewUser(*gormDb)

			result, err := repository.Get(tt.id)

			tt.asserts(t, result, err)
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockedBehavior func(t *testing.T, mock sqlmock.Sqlmock, id string)
		asserts        func(t *testing.T, user *models.User, err error)
	}{
		{
			name: "Test with valid id should delete user",
			id:   "1",
			mockedBehavior: func(t *testing.T, mock sqlmock.Sqlmock, id string) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			asserts: func(t *testing.T, user *models.User, err error) {
				assert.NotNil(t, user)
				assert.Nil(t, err)
			},
		},
		{
			name: "Test with a valid id to delete should return error from db",
			id:   "1",
			mockedBehavior: func(t *testing.T, mock sqlmock.Sqlmock, id string) {

				mock.ExpectQuery(regexp.QuoteMeta("UPDATE `users` SET `deleted_at`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")).
					WillReturnError(errors.New("error al intentar eliminar el usuario"))
				mock.ExpectCommit()
			},
			asserts: func(t *testing.T, user *models.User, err error) {
				assert.Nil(t, user)
				assert.NotNil(t, err)
				assert.Equal(t, "error al intentar eliminar el usuario", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDb, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})

			if err != nil {
				t.Error(err.Error())
			}
			gormDb.Debug()

			tt.mockedBehavior(t, mock, tt.id)

			repository := NewUser(*gormDb)

			result, err := repository.Delete(tt.id)

			tt.asserts(t, result, err)
		})
	}

}

func TestUserRepository_Update(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now()
	tests := []struct {
		name           string
		updateRequest  *models.User
		mockedBehavior func(t *testing.T, mock sqlmock.Sqlmock)
		asserts        func(t *testing.T, user *models.User, err error)
	}{
		{
			name: "Test with valid data should update user",
			updateRequest: &models.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				FirstName: "Meze",
				LastName:  "Law",
				Email:     "meze@gmail.com",
				Password:  "password",
			},
			mockedBehavior: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `created_at`=?,`updated_at`=?,`first_name`=?,`last_name`=?,`email`=?,`password`=? WHERE `users`.`deleted_at` IS NULL AND `id` = ?")).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			asserts: func(t *testing.T, user *models.User, err error) {
				assert.NotNil(t, user)
				assert.Nil(t, err)
			},
		},
		{
			name: "Test with a valid data to update should return error from db",
			updateRequest: &models.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				FirstName: "Meze",
				LastName:  "Law",
				Email:     "meze@gmail.com",
				Password:  "password",
			},
			mockedBehavior: func(t *testing.T, mock sqlmock.Sqlmock) {

				mock.ExpectQuery(regexp.QuoteMeta("UPDATE `users` SET `created_at`=?,`updated_at`=?,`first_name`=?,`last_name`=?,`email`=?,`password`=? WHERE `users`.`deleted_at` IS NULL AND `id` = ?")).
					WillReturnError(errors.New("error al intentar eliminar el usuario"))
				mock.ExpectCommit()
			},
			asserts: func(t *testing.T, user *models.User, err error) {
				assert.Nil(t, user)
				assert.NotNil(t, err)
				assert.Equal(t, "error al actualizar el usuario en DB", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDb, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})

			if err != nil {
				t.Error(err.Error())
			}
			gormDb.Debug()

			tt.mockedBehavior(t, mock)

			repository := NewUser(*gormDb)

			result, err := repository.Update(tt.updateRequest)

			tt.asserts(t, result, err)
		})
	}

}

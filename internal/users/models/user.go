package models

import "time"

type User struct {
	//gorm.Model
	Id        int // TODO UUID?
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

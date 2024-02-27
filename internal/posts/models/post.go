package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title         string
	UserOwnerID   int // TODO
	Amount        float64
	Description   string
	WorkerUserID  int
	PostType      string
	Status        string
	PaymentStatus string
}

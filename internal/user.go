package internal

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type Post struct {
	gorm.Model
	Title  string
	Text   string
	UserID int
	User   User
}

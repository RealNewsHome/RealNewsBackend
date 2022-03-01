package internal

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	Posts    []Post
}

//how to define a method in go: (u === go equivalent of 'this' in JS)
func (u User) GetLocation() {

}

type Post struct {
	gorm.Model
	Title   string
	Text    string
	UserID  int
	User    User
	Upvotes int
}

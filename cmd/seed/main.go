package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/spieziocaroline/realnewsbackend/internal"
	// "github.com/jinzhu/gorm"
)

var (
	users = []internal.User{
		{Username: "Paul", Email: "Paul@gmail.com", Password: "test"},
		{Username: "Diane", Email: "Diane@gmail.com", Password: "test"},
		{Username: "JoJo", Email: "JoJo@gmail.com", Password: "test"},
	}

	posts = []internal.Post{
		{Title: "The voting poll lines are long in Brighton Beach today", Text: "The line to vote at PS431 wrapped around the corner. But it moved pretty fast.", UserID: 1},
		{Title: "Milk prices at $4 a gallon", Text: "I paid $4 for a gallon of regular skim milk today at Fairway. Receipt attached.", UserID: 3},
		{Title: "Don't need a coat today", Text: "Temperature is finally back up to 65 Farenheit in Manhattan.", UserID: 2},
	}
)

func main() {
	dsn := "host=localhost port=5432 user=carolinespiezio dbname=realnews sslmode=disable password=123"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&internal.User{})
	db.AutoMigrate(&internal.Post{})

	//for each car, create a car
	for _, user := range users {
		db.Create(&user)
	}

	for _, post := range posts {
		db.Create(&post)
	}

	fmt.Println("hello sir")
}

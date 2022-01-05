package main

import (
	// "encoding/json"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spieziocaroline/realnewsbackend/internal"
)

var db *gorm.DB
var err error

func main() {

	dsn := "host=localhost port=5432 user=carolinespiezio dbname=realnews sslmode=disable password=123"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// db, err = gorm.Open("postgres", "host=localhost port=5432 user=carolinespiezio dbname=realnews sslmode=disable password=123")

	if err != nil {
		panic("failed to connect database")
	}

	//when a user goes to X route, trigger Y function
	router := mux.NewRouter()
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/newuser", CreateUser).Methods("POST")
	// router.HandleFunc("/posts", GetPosts).Methods("GET")
	// router.HandleFunc("/posts/{id}", GetPost).Methods("GET")
	//get posts by user
	//create post
	//delete post(admins only)

	// router.HandleFunc("/cars", GetCars).Methods("GET")
	// router.HandleFunc("/cars/{id}", GetCar).Methods("GET")
	// router.HandleFunc("/drivers/{id}", GetDriver).Methods("GET")
	// router.HandleFunc("/cars/{id}", DeleteCar).Methods("DELETE")

	//kinda middleware .. browser requires u respond to ceratin requests .. says take this router and wrap it w cors , wrap router w cors stuff

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	fmt.Println("server listening on port8080")

	log.Fatal(http.ListenAndServe(":8080", handler))
}

//here is where we spell out what those 'Y' functions actually mean

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []internal.User
	db.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user internal.User
	db.First(&user, params["id"])
	json.NewEncoder(w).Encode(&user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read the body: ", err)
	}

	var user internal.User
	json.Unmarshal(reqBody, &user)
	if e := db.Create(&user).Error; e != nil {
		log.Println("Unable to create new todo")
	}

	fmt.Println("EndPoint activated! Create New User!")
	json.NewEncoder(w).Encode(user)

	// params := mux.Vars(r)
	// fmt.Println("we are hitting", params)
	// user := internal.User{Username: params["username"], Email: params["email"], Password: params["password"]}

	// db.Create(&user)

	// var users []internal.User
	// db.Find(&users)
	// json.NewEncoder(w).Encode(&users)
}

// func Get(w http.ResponseWriter, r *http.Request) {
// 	//returns the route variables from mux
// 	params := mux.Vars(r)

// 	var driver Driver
// 	var cars []Car

// 	//find the first instance of a type driver that has the id
// 	db.First(&driver, params["id"])
// 	//this ... seems to be finding all cars related to the driver ? and all the data on the driver?
// 	db.Model(&driver).Association("Cars").Find(&cars)
// 	driver.Cars = cars
// 	json.NewEncoder(w).Encode(&driver)
// }

// func DeleteCar(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	var car Car
// 	db.First(&car, params["id"])
// 	db.Delete(&car)

// 	var cars []Car
// 	db.Find(&cars)
// 	json.NewEncoder(w).Encode(&cars)
// }

// func main() {
// 	fmt.Println(mypackage.Add(1, 2))

// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/hello", handleHello)

// 	fmt.Println("Running server on port 8080...")
// 	http.ListenAndServe(":8080", mux)
// }

// func handleHello(w http.ResponseWriter, r *http.Request) {
// 	body, err := io.ReadAll(r.Body)
// 	defer r.Body.Close()
// 	if err != nil {
// 		w.Write([]byte(err.Error()))
// 		w.WriteHeader(500)
// 		return
// 	}

// 	msg := fmt.Sprintf("Hello %s!", string(body))
// 	w.Write([]byte(msg))
// }

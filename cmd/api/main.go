package main

import (
	// "encoding/json"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spieziocaroline/realnewsbackend/internal"
	"golang.org/x/crypto/bcrypt"
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
	router.HandleFunc("/posts", GetPosts).Methods("GET")
	router.HandleFunc("/posts/byUser/{user_id}", GetPostsByUser).Methods("GET")
	router.HandleFunc("/post/{id}", GetPostById).Methods("GET")
	router.HandleFunc("/post", CreatePost).Methods("POST")
	router.HandleFunc("/auth/login", GetAuth).Methods("POST")
	// router.HandleFunc("/auth/me", GetMe).Methods("GET")

	//kinda middleware .. browser requires u respond to ceratin requests .. says take this router and wrap it w cors , wrap router w cors stuff

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	fmt.Println("server listening on port8080")

	log.Fatal(http.ListenAndServe(":8080", handler))
}

//here is where we spell out what those 'Y' functions actually mean

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CreateToken(userid uint64) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func Authenticate(email string, password string) string {
	var user internal.User

	db.Where("email = ?", email).First(&user)
	// hashedPassword, err := HashPassword(password)
	// log.Println("password", password)
	// log.Println("hashedP", hashedPassword)

	token, err := CreateToken(uint64(user.ID))
	log.Println(err)
	if password == user.Password {
		return token
		//find a way to send token to getAuth so we can get user auth .. then pass this to front end & store in get context in App
	} else {
		//return an error - incorrect password!
		return "Incorrect password!"
	}
}

//unhash password?
//and get auth ???
func GetAuth(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read the body: ", err)
	}

	var user internal.User
	json.Unmarshal(reqBody, &user)
	var TOKEN = Authenticate(user.Email, user.Password)
	log.Println(TOKEN)
	json.NewEncoder(w).Encode(TOKEN)
}

// func GetMe(w http.ResponseWriter, r *http.Request) {
// 	token := r.Header.authorization
// 	data := jwt.verify(token, process.env.JWT)
// 	var user internal.User
// 	db.First(&user, data.user)
// 	return user
// }

//get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []internal.User
	db.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

//get specific user
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user internal.User
	db.First(&user, params["id"])
	json.NewEncoder(w).Encode(&user)
}

//create a user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read the body: ", err)
	}

	var user internal.User
	json.Unmarshal(reqBody, &user)
	var oldPassword = user.Password
	hashedPassword, err := HashPassword(oldPassword)
	user.Password = hashedPassword

	if e := db.Create(&user).Error; e != nil {
		log.Println("Unable to create new todo")
	}

	fmt.Println("EndPoint activated! Create New User!")
	json.NewEncoder(w).Encode(user)
}

//get all posts
func GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []internal.Post
	db.Find(&posts)
	json.NewEncoder(w).Encode(&posts)
}

//get just posts from a specific user - if their name is clicked
func GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user internal.User
	var posts []internal.Post
	db.First(&user, params["id"])
	db.Model(&user).Association("Posts").Find(&posts)
	json.NewEncoder(w).Encode(&posts)
}

//get a specific post
func GetPostById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var post internal.Post
	db.First(&post, params["id"])
	json.NewEncoder(w).Encode(&post)
}

//write a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read the body: ", err)
	}

	var post internal.Post
	json.Unmarshal(reqBody, &post)
	if err := db.Create(&post).Error; err != nil {
		log.Println("Unable to create new post")
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Println("EndPoint activated! Create New Post!")
	json.NewEncoder(w).Encode(post)
}

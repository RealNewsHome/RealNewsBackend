package main

import (
	// "encoding/json"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	router.HandleFunc("/auth/me", GetMe).Methods("POST")
	router.HandleFunc("/post/{id}", IncreaseUpvote).Methods("PUT")

	//kinda middleware .. browser requires u respond to ceratin requests .. says take this router and wrap it w cors , wrap router w cors stuff

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	handler := c.Handler(router)
	fmt.Println("server listening on port")

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
	token, err := at.SignedString([]byte("jdnfksdmfksd"))
	if err != nil {
		return "", err
	}
	return token, nil
}

func Authenticate(email string, password string) (string, error) {
	var user internal.User

	db.Where("email = ?", email).First(&user)
	// hashedPassword, err := HashPassword(password)
	// log.Println("password", password)
	// log.Println("hashedP", hashedPassword)

	if user.ID == 0 {
		return "", errors.New("User not found!")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	token, err := CreateToken(uint64(user.ID))
	if err != nil {
		return "", err
	}
	return token, nil
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
	token, err := Authenticate(user.Email, user.Password)
	if err != nil {
		//set status code
		w.WriteHeader(401)
		//turn a string into a byte slice
		w.Write([]byte("Incorrect email or password"))
		return
	}
	json.NewEncoder(w).Encode(token)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	// reqBody, err := ioutil.ReadAll(r.Body)
	// log.Println("r", r)
	// log.Println("request", reqBody)
	// log.Println("error", err)
	tokenString := r.Header.Values("authorization")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString[0], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jdnfksdmfksd"), nil
	})

	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}
	log.Println(token)

	userId := claims["user_id"]

	var user internal.User
	db.First(&user, userId)
	json.NewEncoder(w).Encode(&user)

	// log.Println("reqBody", r.Body)
	// data := jwt.verify(token, process.env.JWT)
	// var user internal.User
	// db.First(&user, data.user)
	// return user
}

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
	log.Println("REQBODY", reqBody)
	var oldPassword = user.Password
	hashedPassword, err := HashPassword(oldPassword)
	if err != nil {
		log.Println("Unable to create hash ", err)
	}
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

	db.First(&user, params["user_id"])

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

func IncreaseUpvote(w http.ResponseWriter, r *http.Request) {
	//find post that matches the ID passed in & increase its upvote
	params := mux.Vars(r)
	var post internal.Post
	db.First(&post, params["id"])
	incremented := post.Upvotes + 1
	post.Upvotes = incremented

	db.Save(&post)
	json.NewEncoder(w).Encode(&post)
}

/*
type App struct {
	UserDB UserDB
}

type UserDB interface {
	GetUser(id string) (internal.User, error)
	InsertUser(user internal.User) error
}

type UserDBGorm struct {
	DB *gorm.DB
}

func (u UserDBGorm) GetUser(id string) (internal.User, error) {
	return internal.User{}, errors.New("not implemented yet")
}

func (u UserDBGorm) InsertUser(user internal.User) error {
	return errors.New("not implemented yet")
}

type UserDBLocal struct {
	data map[uint]internal.User
}

func (u UserDBLocal) GetUser(id uint) (internal.User, error) {
	return u.data[id], nil

}

func (u UserDBLocal) InsertUser(user internal.User) error {
	u.data[user.ID] = user
	return nil
}
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
	// "github.com/jinzhu/gorm"
	"github.com/rs/cors"
)

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "oooooomg"
// 	dbname   = "realnews"
// )

type Driver struct {
	gorm.Model
	Name    string
	License string
	Cars    []Car
}

type Car struct {
	gorm.Model
	Year      int
	Make      string
	ModelName string
	DriverID  int
	Driver    Driver
}

var db *gorm.DB
var err error
var (
	drivers = []Driver{
		{Name: "Jimmy Johnson", License: "ABC123"},
		{Name: "Howard Hills", License: "XYZ789"},
		{Name: "Craig Colbin", License: "DEF333"},
	}

	cars = []Car{
		{Year: 2000, Make: "Toyota", ModelName: "Tundra", DriverID: 1},
		{Year: 2001, Make: "Honda", ModelName: "Accord", DriverID: 1},
		{Year: 2002, Make: "Nissan", ModelName: "Sentra", DriverID: 2},
		{Year: 2003, Make: "Ford", ModelName: "F-150", DriverID: 3},
	}
)

func main() {
	router := mux.NewRouter()

	dsn := "host=localhost port=5432 user=carolinespiezio dbname=realnews sslmode=disable password=123"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// db, err = gorm.Open("postgres", "host=localhost port=5432 user=carolinespiezio dbname=realnews sslmode=disable password=123")

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Driver{})
	db.AutoMigrate(&Car{})

	//for each car, create a car
	for _, car := range cars {
		db.Create(&car)
	}

	for _, driver := range drivers {
		db.Create(&driver)
	}

	//when a user goes to X route, trigger Y function
	router.HandleFunc("/cars", GetCars).Methods("GET")
	router.HandleFunc("/cars/{id}", GetCar).Methods("GET")
	router.HandleFunc("/drivers/{id}", GetDriver).Methods("GET")
	router.HandleFunc("/cars/{id}", DeleteCar).Methods("DELETE")

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
func GetCars(w http.ResponseWriter, r *http.Request) {
	var cars []Car
	db.Find(&cars)
	json.NewEncoder(w).Encode(&cars)
}

func GetCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var car Car
	db.First(&car, params["id"])
	json.NewEncoder(w).Encode(&car)
}

func GetDriver(w http.ResponseWriter, r *http.Request) {
	//returns the route variables from mux
	params := mux.Vars(r)

	var driver Driver
	var cars []Car

	//find the first instance of a type driver that has the id
	db.First(&driver, params["id"])
	//this ... seems to be finding all cars related to the driver ? and all the data on the driver?
	db.Model(&driver).Association("Cars").Find(&cars)
	driver.Cars = cars
	json.NewEncoder(w).Encode(&driver)
}

func DeleteCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var car Car
	db.First(&car, params["id"])
	db.Delete(&car)

	var cars []Car
	db.Find(&cars)
	json.NewEncoder(w).Encode(&cars)
}

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

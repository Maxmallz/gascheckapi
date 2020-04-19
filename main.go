package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var cars []Car

func main() {
	//initializePG()
	appendCars()
	initializeRouter()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func appendCars() {
	cars = append(cars,
		Car{Make: "Fiat", Model: "500X", Year: "2020"},
		Car{Make: "Ford", Model: "Focus", Year: "2016"})
}

func dbHealth(w http.ResponseWriter, r *http.Request) {
	gormDb, err := gorm.Open("postgres", "user=postgres password=g@schekeR! dbname=GasCheck sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer gormDb.Close()

	db := gormDb.DB()

	err = db.Ping()

	if err != nil {
		json.NewEncoder(w).Encode("db unhealty")
		return
	}

	json.NewEncoder(w).Encode("db healthy")
}

func initializeRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/cars", getCars).Methods("GET")
	router.HandleFunc("/dbHealth", dbHealth).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getCars(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cars)
}

type Car struct {
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  string `json:"year"`
}

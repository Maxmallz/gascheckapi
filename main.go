package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
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

func gormHealth(w http.ResponseWriter, r *http.Request) {
	gormDb, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		json.NewEncoder(w).Encode("go orm error")
		return
	}

	defer gormDb.Close()

	db := gormDb.DB()

	err = db.Ping()

	if err != nil {
		json.NewEncoder(w).Encode("go orm error")
		return
	}

	json.NewEncoder(w).Encode("go orm success")
}

func initializeRouter() {
	router := mux.NewRouter()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	print("app is listening on port " + port)

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/cars", getCars).Methods("GET")
	router.HandleFunc("/dbhealth", dbHealth).Methods("GET")
	router.HandleFunc("/gorm", gormHealth).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getCars(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cars)
}

func index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("app is running")
}

func dbHealth(w http.ResponseWriter, r *http.Request) {
	dbDev, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	print(os.Getenv("DATABASE_URL"))

	var dbCons = [2]string{}

	if err != nil {
		dbCons[0] = "hobby-dev connection failed"
		log.Fatalf("Error opening database: %q", err)
	}

	defer dbDev.Close()

	err = dbDev.Ping()

	if err != nil {
		dbCons[0] = "hobby-dev not reacheable"
		log.Fatalf("Error opening database: %q", err)
	} else {
		dbCons[0] = "hobby-dev connected"
	}

	dbBasic, err := sql.Open("postgres", os.Getenv("HEROKU_POSTGRESQL_COBALT_URL"))

	print(os.Getenv("HEROKU_POSTGRESQL_COBALT_URL"))

	if err != nil {
		dbCons[1] = "hobby-basic connection failed"
		log.Fatalf("Error opening database: %q", err)
	}

	defer dbBasic.Close()
	err = dbBasic.Ping()

	if err != nil {
		dbCons[1] = "hobby-basic not reacheable"
		log.Fatalf("Error opening database: %q", err)
	} else {
		dbCons[1] = "hobby-basic connected"
	}

	json.NewEncoder(w).Encode(dbCons)
}

type Car struct {
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  string `json:"year"`
}

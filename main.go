package main

import (
	"database/sql"
	"encoding/json"
	"gascheckapi/config"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

func main() {
	config.ConfigureEnv()
	initializeRouter()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
		port = "8000"
	}

	print("app is listening on port " + port + "\n")

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/dbhealth", dbHealth).Methods("GET")
	router.HandleFunc("/api", getApi).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("app is running")
}

func getApi(w http.ResponseWriter, r *http.Request) {
	var apis = []string{
		"api: /dbhealth, desc: list avaialable api's",
		"api: /, desc: check app state",
		"api: /dbhealth, desc: check db health"}

	json.NewEncoder(w).Encode(apis)
}

func dbHealth(w http.ResponseWriter, r *http.Request) {
	con := config.GetConnStr()
	db, err := sql.Open("postgres", con)

	if err != nil {
		log.Fatal(err)
		json.NewEncoder(w).Encode("error opening database")
		return
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		json.NewEncoder(w).Encode("ping failed")
		return
	}

	json.NewEncoder(w).Encode("aws rds db connected")
}

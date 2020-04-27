package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

var cars []Car

func main() {
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
		Car{Make: "Ford", Model: "Focus", Year: "2016"},
		Car{Make: "Nissan", Model: "Versa", Year: "2020"})
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

	print("app is listening on port " + port)

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/cars", getCars).Methods("GET")
	router.HandleFunc("/dbhealth", dbHealth).Methods("GET")
	router.HandleFunc("/gorm", gormHealth).Methods("GET")
	router.HandleFunc("/tick", tick).Methods("GET")
	//router.HandleFunc("/envVar", getEnvVar).Methods("GET")
	router.HandleFunc("/env", getEnv).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getCars(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cars)
}

func index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("app is running")
}

func dbHealth(w http.ResponseWriter, r *http.Request) {
	// env := config.GetEnv()

	// c, err := config.GetConfig(env)

	// if err != nil {
	// 	logFatal(err)
	// }
	// con := config.GetConnStr(c)
	// db, err := sql.Open("postgres", con)

	// if err != nil {
	// 	log.Fatal(err)
	// 	json.NewEncoder(w).Encode("error opening database")
	// 	return
	// }

	// defer db.Close()

	// err = db.Ping()

	// if err != nil {
	// 	json.NewEncoder(w).Encode("ping failed")
	// 	return
	// }

	// json.NewEncoder(w).Encode("aws rds db connected")
}

func tick(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	print(os.Getenv("DATABASE_URL"))

	if err != nil {
		json.NewEncoder(w).Encode("error opening database")
		return
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("error creating database table: %q", err))
		return
	}

	if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("error incrementing tick: %q", err))
		return
	}

	rows, err := db.Query("SELECT tick FROM ticks")
	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("error reading ticks: %q", err))
		return
	}

	defer rows.Close()

	for rows.Next() {
		var tick time.Time
		if err := rows.Scan(&tick); err != nil {
			json.NewEncoder(w).Encode("error scanning ticks")
			return
		}

		json.NewEncoder(w).Encode(fmt.Sprintf("Read from DB: %s\n", tick.String()))
	}
}

// func getEnvVar(w http.ResponseWriter, r *http.Request) {
// 	env := config.GetEnv()
// 	json.NewEncoder(w).Encode(env)
// }

func getEnv(w http.ResponseWriter, r *http.Request) {
	arr := []string{os.Getenv("APP_ENV"), os.Getenv("PASSWORD"), os.Getenv("SSLMODE"), os.Getenv("USERNAME")}
	json.NewEncoder(w).Encode(arr)
}

type Car struct {
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  string `json:"year"`
}

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	Environment string
	DbHost      string
	DbPort      int
	Database    string
	Username    string
	Password    string
	SSLMode     string
}

var envs = []string{"local", "production"}
var c *Config

// Processes the command line arguments for application startup and returns the requested environment.
// Use 'env=' to specify target environment. 'env=local' is default
func getEnv() string {
	e := "local"
	args := os.Args

	if len(args) < 1 {
		return e
	}

	sort.Strings(args)

	i, found := find(args, "env=")

	if !found || i == -1 {
		return e
	}

	tempEnv := strings.SplitAfter(args[i], "env=")[1]

	i, found = find(envs, tempEnv)

	if found {
		e = tempEnv
	}

	return e
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if strings.Contains(item, val) {
			return i, true
		}
	}

	return -1, false
}

func buiildConfig(env string) (*Config, error) {
	b, err := ioutil.ReadFile("config/config.json")

	if err != nil {
		log.Fatal(err)
	}

	c := []Config{}

	if json.Valid(b) {
		err = json.Unmarshal(b, &c)

		if err != nil {
			return nil, errors.New("Failed to unmarshal config file")
		}

		for _, config := range c {
			if config.Environment == env {
				return &config, nil
			}
		}
	}

	return nil, errors.New("Failed to load config file")
}

func GetConnStr() string {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.DbHost, c.DbPort, c.Username, c.Password, c.Database, c.SSLMode)
	return conn
}

func ConfigureEnv() {
	if os.Getenv("APP_ENV") == "production" {
		port, err := strconv.Atoi(os.Getenv("DB_PORT"))

		if err != nil {
			panic(err)
		}

		c = &Config{
			Environment: os.Getenv("APP_ENV"),
			Database:    os.Getenv("DATABASE"),
			DbHost:      os.Getenv("DB_HOST"),
			DbPort:      port,
			Username:    os.Getenv("USERNAME"),
			Password:    os.Getenv("PASSWORD"),
			SSLMode:     os.Getenv("SSLMODE")}

		return
	}

	env := getEnv()
	conf, err := buiildConfig(env)

	if err != nil {
		panic(err)
	}

	c = conf
}

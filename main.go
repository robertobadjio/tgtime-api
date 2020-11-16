package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"officetime-api/app/aggregator"
	"officetime-api/app/config"
	"officetime-api/app/dao"
	"officetime-api/app/model"
	"time"
)

var db *sql.DB

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "The config name param")
	flag.Parse()

	if configPath == "" {
		fmt.Println("Param 'config' must be set")
		return
	}
	if err := config.LoadConfig(configPath); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	db = getDB()
	dao.Db = db
	model.Db = db
	aggregator.Db = db
	go every12Day()

	fmt.Println("Setting up server, enabling CORS...")
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedMethods: []string{"GET", "POST", "PATCH"}, // Allowing only get, just an example
	})

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api-service/time/{id}/day/{date}", dao.GetTimeDayAll).Methods("GET")
	router.HandleFunc("/api-service/time/{id}/period/{period}", dao.GetTimeByPeriod).Methods("GET")
	router.HandleFunc("/api-service/time", dao.CreateTime).Methods("POST")

	router.HandleFunc("/api-service/period", dao.GetAllPeriods).Methods("GET")
	router.HandleFunc("/api-service/period/{id}", dao.GetPeriod).Methods("GET")
	router.HandleFunc("/api-service/period", dao.CreatePeriod).Methods("POST")
	router.HandleFunc("/api-service/period/{id}", dao.UpdatePeriod).Methods("PATCH")
	router.HandleFunc("/api-service/period/{id}", dao.DeletePeriod).Methods("DELETE")

	router.HandleFunc("/api-service/user", dao.GetAllUsers).Methods("GET")
	router.HandleFunc("/api-service/user/{id}", dao.GetUser).Methods("GET")
	router.HandleFunc("/api-service/user/{id}", dao.UpdateUser).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":8080", c.Handler(router)))
}

func every12Day() {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
	d := n.Sub(t)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	for {
		time.Sleep(d)
		d = 24 * time.Hour
		aggregator.AggregateTime()
	}
}

func getDB() *sql.DB {
	pgConString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		config.Config.HostName, config.Config.HostPort, config.Config.UserName, config.Config.Password, config.Config.DataBaseName, config.Config.SslMode)

	db, err := sql.Open("postgres", pgConString)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}

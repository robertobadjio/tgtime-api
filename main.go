package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"officetime-api/app/config"
	"officetime-api/app/dao"
	"github.com/gorilla/mux"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "The config name param")
	flag.Parse()
	fmt.Println()

	if configPath == "" {
		fmt.Println("Param 'config' must be set")
		return
	}

	if err := config.LoadConfig(configPath); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/time/{id}/days/{date}", dao.GetTimeDayAll).Methods("GET")
	router.HandleFunc("/time", dao.CreateTime).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
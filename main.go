package main

import (
	"database/sql"
	"flag"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"log"
	"net/http"
	"officetime-api/app/aggregator"
	"officetime-api/app/config"
	"officetime-api/app/dao"
	"officetime-api/app/model"
	"officetime-api/app/service"
	"strconv"
	"time"
)

var db *sql.DB

// Global secret key
var mySigningKey = []byte(config.Config.AuthSigningKey)

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
		AllowedOrigins: []string{"*"},                    // All origins
		AllowedMethods: []string{"GET", "POST", "PATCH"}, // Allowing only get, just an example
	})

	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	router.HandleFunc("/api-service/login", dao.Login).Methods("POST")
	router.HandleFunc("/api-service/token/refresh", dao.Refresh).Methods("POST")
	router.HandleFunc("/api-service/logout", dao.Logout).Methods("POST")

	router.Handle("/metrics", promhttp.Handler())

	router.Handle("/api-service/time/{id}/day/{date}", isAuthorized(dao.GetTimeDayAll)).Methods("GET")
	router.Handle("/api-service/time/{id}/period/{period}", isAuthorized(dao.GetTimeByPeriod)).Methods("GET")
	router.Handle("/api-service/time", isAuthorized(dao.CreateTime)).Methods("POST")

	router.Handle("/api-service/period", isAuthorized(dao.GetAllPeriods)).Methods("GET")
	router.Handle("/api-service/period/{id}", isAuthorized(dao.GetPeriod)).Methods("GET")
	router.Handle("/api-service/period", isAuthorized(dao.CreatePeriod)).Methods("POST")
	router.Handle("/api-service/period/{id}", isAuthorized(dao.UpdatePeriod)).Methods("PATCH")
	router.Handle("/api-service/period/{id}", isAuthorized(dao.DeletePeriod)).Methods("DELETE")

	router.Handle("/api-service/user", isAuthorized(dao.GetAllUsers)).Methods("GET")
	router.Handle("/api-service/user/{id}", isAuthorized(dao.GetUser)).Methods("GET")
	router.Handle("/api-service/user/{id}", isAuthorized(dao.UpdateUser)).Methods("PATCH")
	router.Handle("/api-service/user", isAuthorized(dao.CreateUser)).Methods("POST")
	router.Handle("/api-service/user/{id}", isAuthorized(dao.DeleteUser)).Methods("DELETE")

	router.Handle("/api-service/department/{id}", isAuthorized(dao.GetDepartment)).Methods("GET")
	router.Handle("/api-service/department", isAuthorized(dao.GetAllDepartments)).Methods("GET")
	router.Handle("/api-service/department", isAuthorized(dao.CreateDepartment)).Methods("POST")
	router.Handle("/api-service/department/{id}", isAuthorized(dao.UpdateDepartment)).Methods("PATCH")
	router.Handle("/api-service/department/{id}", isAuthorized(dao.DeleteDepartment)).Methods("DELETE")

	router.Handle("/api-service/router/{id}", isAuthorized(dao.GetRouter)).Methods("GET")
	router.Handle("/api-service/router", isAuthorized(dao.GetAllRouters)).Methods("GET")
	router.Handle("/api-service/router", isAuthorized(dao.CreateRouter)).Methods("POST")
	router.Handle("/api-service/router/{id}", isAuthorized(dao.UpdateRouter)).Methods("PATCH")
	router.Handle("/api-service/router/{id}", isAuthorized(dao.DeleteRouter)).Methods("DELETE")

	router.Handle("/api-service/stat/periods-and-routers", isAuthorized(dao.GetStatByPeriodsAndRouters)).Methods("GET")
	router.Handle("/api-service/stat/departments/{date}", isAuthorized(dao.GetAllTimesDepartmentsByDate)).Methods("GET")

	router.Handle("/api-service/weekend", isAuthorized(dao.GetWeekend)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", c.Handler(router)))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func every12Day() {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 0, 1, 0, 0, t.Location())
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

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				if "Token is expired" == err.Error() {
					w.WriteHeader(http.StatusUnauthorized)
				}
				fmt.Fprintf(w, err.Error())
				return
			}

			au, err := service.ExtractTokenMetadata(r)
			if err != nil {
				if "Token is expired" == err.Error() {
					w.WriteHeader(http.StatusUnauthorized)
				}
				fmt.Fprintf(w, err.Error())
				return
			}
			if au == nil {
				fmt.Fprintf(w, err.Error())
				if "Token is expired" == err.Error() {
					w.WriteHeader(http.StatusUnauthorized)
				}
				return
			}

			id := mux.Vars(r)["id"]
			userId, _ := strconv.Atoi(id) // TODO: если обрабатывать ошибку, апи падает

			if au.UserId != uint64(userId) && !au.IsAdmin() {
				fmt.Fprintf(w, "Access denied")
				return
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

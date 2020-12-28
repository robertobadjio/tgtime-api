package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"officetime-api/app/aggregator"
	"officetime-api/app/config"
	"officetime-api/app/dao"
	"officetime-api/app/model"
	"strconv"
	"strings"
	"time"
)

type authData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccessDetails struct {
	UserId uint64
	Role   string
}

func (ad *AccessDetails) isAdmin() bool {
	return "Admin" == ad.Role // TODO: const
}

var db *sql.DB

// Global secret key
var mySigningKey = []byte("vtlcgjgek") // TODO: ключ в конфиг

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update")
	}

	var data authData
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		panic(err)
	}

	if !dao.CheckAuth(data.Email, data.Password) {
		json.NewEncoder(w).Encode("Invalid user")
		return
	}
	// Создаем новый токен
	token := jwt.New(jwt.SigningMethodHS256)

	user := model.GetUserByEmail(data.Email)
	claims := token.Claims.(jwt.MapClaims)
	// Устанавливаем набор параметров для токена
	claims["authorized"] = true
	claims["userId"] = user.Id
	claims["userName"] = user.Name
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix() // TODO: время в конфиг
	claims["role"] = user.Role // TODO: костыль RBAC?

	// Подписываем токен нашим секретным ключем
	tokenString, _ := token.SignedString(mySigningKey)

	// Отдаем токен клиенту
	w.Write([]byte(tokenString))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	_, err := ExtractTokenMetadata(r) // TODO: au

	if err != nil {
		w.Write([]byte("Successfully logged out"))
		return
	}

	// TODO: сделать разлогин через BlackWhite lists
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["userId"]), 10, 64)
		if err != nil {
			return nil, err
		}
		role := fmt.Sprintf("%s", claims["role"])
		return &AccessDetails{
			UserId: userId,
			Role: role,
		}, nil
	}
	return nil, err
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Token")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 1 {
		return strArr[0]
	}
	return ""
}

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

	router.HandleFunc("/api-service/login", GetTokenHandler).Methods("POST")
	router.HandleFunc("/api-service/logout", Logout).Methods("POST")

	router.Handle("/api-service/time/{id}/day/{date}", isAuthorized(dao.GetTimeDayAll)).Methods("GET")
	router.Handle("/api-service/time/{id}/period/{period}", isAuthorized(dao.GetTimeByPeriod)).Methods("GET")
	router.HandleFunc("/api-service/time", dao.CreateTime).Methods("POST")

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

	router.Handle("/api-service/department", isAuthorized(dao.GetAllDepartments)).Methods("GET")
	router.Handle("/api-service/department/{id}", isAuthorized(dao.UpdateDepartment)).Methods("PATCH")
	router.Handle("/api-service/department/{id}", isAuthorized(dao.DeleteDepartment)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", c.Handler(router)))
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
				fmt.Fprintf(w, err.Error())
			}

			au, err := ExtractTokenMetadata(r)
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			if au == nil {
				fmt.Fprintf(w, err.Error())
				return
			}

			id := mux.Vars(r)["id"]
			userId, _ := strconv.Atoi(id)
			// TODO: обработка ошибок

			if au.UserId != uint64(userId) && !au.isAdmin() {
				fmt.Fprintf(w, "Access denied")
				return
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"officetime-api/app/service"
	"officetime-api/internal/config"
	"officetime-api/internal/db"
	"officetime-api/internal/model/user/adapter"
	userApp "officetime-api/internal/model/user/app"
	"officetime-api/internal/model/user/app/query"
	"strconv"
	"time"
)

type RefreshToken struct {
	Token string `json:"refresh_token"`
}

type authData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenDetails struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenExpires  int64  `json:"access_token_expires"`
	RefreshTokenExpires int64  `json:"refresh_token_expires"`
}

func Logout(w http.ResponseWriter, r *http.Request) {
	_, err := service.ExtractTokenMetadata(r) // TODO: au
	if err != nil {
		w.Write([]byte("Successfully logged out"))
		return
	}

	// TODO: сделать разлогин через BlackWhite lists
}

func Login(w http.ResponseWriter, r *http.Request) {
	cfg := config.New()

	td := &TokenDetails{}
	td.AccessTokenExpires = time.Now().Add(time.Minute * time.Duration(cfg.AuthAccessTokenExpires)).Unix()
	td.RefreshTokenExpires = time.Now().Add(time.Hour * time.Duration(cfg.AuthRefreshTokenExpires)).Unix()

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update")
	}

	var data authData
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		panic(err)
	}

	uApp := buildUserApp()
	qr := query.GetUserByEmail{Email: data.Email}
	ctx := context.TODO()
	user, err := uApp.Queries.GetUserByEmail.Handle(ctx, qr)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	qr2 := query.GetUserPasswordHashByEmail{Email: user.Email}
	// TODO: Hadle error
	userPasswordHash, _ := uApp.Queries.GetUserPasswordHashByEmail.Handle(ctx, qr2) // TODO: Убрать
	if !service.CheckAuth(userPasswordHash, data.Password) {
		fmt.Fprintf(w, "Wrong password")
		return
	}

	json.NewEncoder(w).Encode(service.CreateTokenPair(user))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	cfg := config.New()

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var refreshToken RefreshToken
	err = json.Unmarshal(reqBody, &refreshToken)
	if err != nil {
		panic(err)
	}

	// Verify the token
	token, err := jwt.Parse(refreshToken.Token, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.AuthRefreshKey), nil
	})

	w.Header().Set("Content-Type", "application/json")
	// If there is an error, the token must have expired
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) // "Refresh token expired"
		return
	}

	// Is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["userId"]), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity) // TODO: вернуть "Error occurred"
			fmt.Println("Error occurred")
			//w.Write([]byte("Error occurred"))
			return
		}

		qr := query.GetUser{UserId: int(userId)}
		ctx := context.TODO()
		app := buildUserApp()
		u, _ := app.Queries.GetUser.Handle(ctx, qr) // TODO: Handle error
		// Create new pairs of refresh and access tokens

		json.NewEncoder(w).Encode(service.CreateTokenPair(u)) // TODO: обработка ошибки, если пользователь не найден
	}
}

func buildUserApp() userApp.Application {
	userRepository := adapter.NewPgUserRepository(db.GetDB())
	return userApp.Application{
		Queries: userApp.Queries{
			GetUser:                    query.NewGetUserHandler(userRepository),
			GetUserByEmail:             query.NewGetUserByEmailHandler(userRepository),
			GetUserPasswordHashByEmail: query.NewGetUserPasswordHashByEmailHandler(userRepository),
		},
	}
}

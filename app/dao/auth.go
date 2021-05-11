package dao

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"officetime-api/app/config"
	"officetime-api/app/model"
	"officetime-api/app/service"
	"time"
)

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
	td := &TokenDetails{}
	td.AccessTokenExpires = time.Now().Add(time.Minute * time.Duration(config.Config.AuthAccessTokenExpires)).Unix()
	td.RefreshTokenExpires = time.Now().Add(time.Hour * time.Duration(config.Config.AuthRefreshTokenExpires)).Unix()

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the mac address and seconds only in order to update")
	}

	var data authData
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		panic(err)
	}

	user, err := model.GetUserByEmail(data.Email)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	userPasswordHash := model.GetUserPasswordHashByEmail(user.Email) // TODO: Убрать
	if !service.CheckAuth(userPasswordHash, data.Password) {
		fmt.Fprintf(w, "Wrong password")
		return
	}

	json.NewEncoder(w).Encode(service.CreateTokenPair(user))
}

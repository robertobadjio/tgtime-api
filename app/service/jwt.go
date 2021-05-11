package service

import (
	"github.com/dgrijalva/jwt-go"
	"officetime-api/app/config"
	"officetime-api/app/model"
	"time"
)

// Global secret key
var mySigningKey = []byte(config.Config.AuthSigningKey)
var refreshSecretKey = []byte(config.Config.AuthRefreshKey)

type TokenDetails struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenExpires  int64  `json:"access_token_expires"`
	RefreshTokenExpires int64  `json:"refresh_token_expires"`
}

func CreateTokenPair(user *model.User) *TokenDetails {
	td := &TokenDetails{}
	td.AccessTokenExpires = time.Now().Add(time.Minute * time.Duration(config.Config.AuthAccessTokenExpires)).Unix()
	td.RefreshTokenExpires = time.Now().Add(time.Hour * time.Duration(config.Config.AuthRefreshTokenExpires)).Unix()

	// Создаем новый токен
	token := jwt.New(jwt.SigningMethodHS256)

	accessTokenClaims := token.Claims.(jwt.MapClaims)
	// Устанавливаем набор параметров для токена
	accessTokenClaims["authorized"] = true
	accessTokenClaims["userId"] = user.Id
	accessTokenClaims["userFirstname"] = user.Name
	accessTokenClaims["userSurname"] = user.Surname
	accessTokenClaims["userLastname"] = user.Lastname
	accessTokenClaims["userEmail"] = user.Email
	accessTokenClaims["userBirthDate"] = user.BirthDate
	accessTokenClaims["exp"] = td.AccessTokenExpires
	accessTokenClaims["role"] = user.Role // TODO: костыль, RBAC?
	accessTokenClaims["department"] = user.Department
	accessTokenClaims["position"] = user.Position

	// Подписываем токен нашим секретным ключем
	var err error
	td.AccessToken, err = token.SignedString(mySigningKey)
	if err != nil {
		panic(err)
	}

	// Creating Refresh Token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["userId"] = user.Id
	refreshTokenClaims["exp"] = td.RefreshTokenExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	td.RefreshToken, err = refreshToken.SignedString(refreshSecretKey)
	if err != nil {
		panic(err)
	}

	return td
}

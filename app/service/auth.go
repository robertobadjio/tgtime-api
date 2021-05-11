package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"officetime-api/app/config"
	"officetime-api/app/model"
	"strconv"
	"strings"
	"time"
)

// Global secret key
var mySigningKey = []byte(config.Config.AuthSigningKey)
var refreshSecretKey = []byte(config.Config.AuthRefreshKey)

type AccessDetails struct {
	UserId uint64
	Role   string
}

func (ad *AccessDetails) IsAdmin() bool {
	return "Admin" == ad.Role // TODO: const
}

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

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)

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

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Token")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 1 {
		return strArr[0]
	}
	return ""
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
			Role:   role,
		}, nil
	}
	return nil, err
}
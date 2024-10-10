package service

/*
import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/robertobadjio/tgtime-api/internal/config"
	"github.com/robertobadjio/tgtime-api/internal/model/user/domain/user"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

func CreateTokenPair(u *user.User) *TokenDetails {
	cfg := config.New()
	td := &TokenDetails{}
	td.AccessTokenExpires = time.Now().Add(time.Minute * time.Duration(cfg.AuthAccessTokenExpires)).Unix()
	td.RefreshTokenExpires = time.Now().Add(time.Hour * time.Duration(cfg.AuthRefreshTokenExpires)).Unix()

	// Создаем новый токен
	token := jwt.New(jwt.SigningMethodHS256)

	accessTokenClaims := token.Claims.(jwt.MapClaims)
	// Устанавливаем набор параметров для токена
	accessTokenClaims["authorized"] = true
	accessTokenClaims["userId"] = u.Id
	accessTokenClaims["userFirstname"] = u.Name
	accessTokenClaims["userSurname"] = u.Surname
	accessTokenClaims["userLastname"] = u.Lastname
	accessTokenClaims["userEmail"] = u.Email
	accessTokenClaims["userBirthDate"] = u.BirthDate
	accessTokenClaims["exp"] = td.AccessTokenExpires
	accessTokenClaims["role"] = u.Role // TODO: костыль, RBAC?
	accessTokenClaims["department"] = u.Department
	accessTokenClaims["position"] = u.Position

	// Подписываем токен нашим секретным ключем
	var err error
	td.AccessToken, err = token.SignedString([]byte(cfg.AuthSigningKey))
	if err != nil {
		panic(err)
	}

	// Creating Refresh Token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["userId"] = u.Id
	refreshTokenClaims["exp"] = td.RefreshTokenExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	td.RefreshToken, err = refreshToken.SignedString([]byte(cfg.AuthRefreshKey))
	if err != nil {
		panic(err)
	}

	return td
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	cfg := config.New()
	tokenString := extractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.AuthSigningKey), nil
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
*/

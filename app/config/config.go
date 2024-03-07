package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"strconv"
)

const projectDirName = "officetime-api"

type Config struct {
	DataBaseHost            string
	DataBasePort            int
	DataBaseName            string
	DataBaseUser            string
	DataBasePassword        string
	DataBaseSslMode         string
	AuthSigningKey          string
	AuthRefreshKey          string
	AuthAccessTokenExpires  int
	AuthRefreshTokenExpires int
}

func init() {
	loadEnv()
}

func New() *Config {
	return &Config{
		DataBaseHost:            getEnv("DATABASE_HOST", ""),
		DataBasePort:            getEnvInt("DATABASE_PORT", 0),
		DataBaseName:            getEnv("DATABASE_NAME", ""),
		DataBaseUser:            getEnv("DATABASE_USER", ""),
		DataBasePassword:        getEnv("DATABASE_PASSWORD", ""),
		DataBaseSslMode:         getEnv("DATABASE_SSL_MODE", ""),
		AuthSigningKey:          getEnv("AUTH_SIGNING_KEY", ""),
		AuthRefreshKey:          getEnv("AUTH_REFRESH_KEY", ""),
		AuthAccessTokenExpires:  getEnvInt("AUTH_ACCESS_TOKEN_EXPIRES", 0),
		AuthRefreshTokenExpires: getEnvInt("AUTH_REFRESH_TOKEN_EXPIRES", 0),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func loadEnv() {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Problem loading .env file")
	}
}

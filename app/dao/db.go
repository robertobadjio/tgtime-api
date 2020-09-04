package dao

import (
	"database/sql"
	"fmt"
	"officetime-api/app/config"
)

func GetDB() *sql.DB {
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
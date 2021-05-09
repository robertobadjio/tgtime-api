package service

import (
	"log"
	"net/smtp"
	"officetime-api/app/config"
)

func SendEmail(name, email, password string) {
	// TODO: В параметры
	auth := smtp.PlainAuth("", "info@officetime.tech", "E61kTuq7", "smtp.timeweb.ru")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Регистрация на OfficeTime\r\n" +
		"Добро пожаловать " + name + "! Ваш пароль: " + password + "\r\n\r\n" +
		"Для начала работы напишите '/start' телеграм-боту @" + config.Config.TelegramBot)
	err := smtp.SendMail("smtp.timeweb.ru:25", auth, "info@officetime.tech", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

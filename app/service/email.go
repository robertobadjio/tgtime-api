package service

import (
	"log"
	"net/smtp"
)

func SendEmail(name, email, password string) {
	// TODO: В параметры
	auth := smtp.PlainAuth("", "info@tgtime.ru", "E61kTuq7", "smtp.timeweb.ru")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Регистрация на TgTime\r\n" +
		"Добро пожаловать " + name + "! Ваш пароль: " + password + "\r\n\r\n" +
		"Для начала работы напишите '/start' телеграм-боту @TgTimeTech")
	err := smtp.SendMail("smtp.timeweb.ru:25", auth, "info@tgtime.ru", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

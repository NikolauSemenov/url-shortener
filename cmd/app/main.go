package main

import (
	"log"
	"url-shortener/internal/app"
)

// @title Url-shortener
// @version 1.0.0
// @description Сервис для сокращения ссылок
func main() {
	serv, err := app.NewApp()
	if err != nil {
		log.Fatal("Ошибка настройки сервера: ", err)
	}
	if err = serv.Run(); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}

}

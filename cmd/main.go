package main

import (
	db2 "go_final_project/internal/db"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("Файл хранения переменных окружения не найден")
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	db := db2.SetupDB()
	defer db.Close()

	// Установка директории для статических файлов
	webDir := "./web"

	// Настройка обработчика файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен и слушает порт %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера", err)

	}
}

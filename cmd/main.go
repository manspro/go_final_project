package main

import (
	"fmt"
	"go_final_project/internal/db"
	"go_final_project/internal/handler"
	"log"
	"net/http"
	"os"
	"time"

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

	db := db.SetupDB()
	defer db.Close()

	// Установка директории для статических файлов
	webDir := "./web"

	// Настройка обработчика файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		nowStr := r.FormValue("now")
		dateStr := r.FormValue("date")
		repeat := r.FormValue("repeat")

		now, err := time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Невалидная 'now' дата: %s", nowStr), http.StatusBadRequest)
			return
		}

		nextDate, err := handler.NextDate(now, dateStr, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte(nextDate))
	})

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handler.TaskHandler(db, w, r)
	})

	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		handler.GetTasks(db, w, r)
	})

	http.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		handler.MarkTaskDone(db, w, r)
	})
	log.Printf("Сервер запущен и слушает порт %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера", err)

	}
}

package handler

import (
	"database/sql"
	"encoding/json"
	"go_final_project/internal/model"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func AddTask(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var task model.Task
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(1212121212, task)

	if task.Title == "" {
		http.Error(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	now := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = now
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, "Дата представлена в неверном формате", http.StatusBadRequest)
			return
		}
		if parsedDate.Before(time.Now()) {
			if task.Repeat != "" {
				// Вызов функции NextDate для корректировки даты
				nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					http.Error(w, "Ошибка при вычислении даты по правилу повторения: "+err.Error(), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			} else {
				task.Date = now
			}
		}
	}

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, "Ошибка добавления задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Ошибка получения ID задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

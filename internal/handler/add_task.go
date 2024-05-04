package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/internal/model"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func AddTask(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	log.Println("Received body:", string(body))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var task model.Task
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&task); err != nil {
		sendJsonError(w, fmt.Sprintf("Error parsing JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		sendJsonError(w, "title required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	todayStr := now.Format("20060102")
	if task.Date == "" {
		task.Date = todayStr
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			sendJsonError(w, "invalid date", http.StatusBadRequest)
			return
		}
		if parsedDate.Before(now) {
			if task.Repeat != "" {
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					sendJsonError(w, fmt.Sprintf("repeat error: %s", err.Error()), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			} else {
				task.Date = todayStr
			}
		}
	}

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		sendJsonError(w, fmt.Sprintf("error add task: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		sendJsonError(w, fmt.Sprintf("error get task by ID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func sendJsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

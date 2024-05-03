package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/internal/model"
	"net/http"
	"time"
)

// GetTasks обрабатывает GET запросы на /api/tasks
func GetTasks(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	var rows *sql.Rows
	var err error

	query := "SELECT id, date, title, comment, repeat FROM scheduler"
	if search != "" {
		date, err := time.Parse("2.01.2006", search)
		if err == nil {
			query += " WHERE date = ?"
			rows, err = db.Query(query, date.Format("20060102"))
		} else {
			searchPattern := "%" + search + "%"
			query += " WHERE title LIKE ? OR comment LIKE ?"
			rows, err = db.Query(query, searchPattern, searchPattern)
		}
	} else {
		rows, err = db.Query(query)
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка базы данных: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []model.GetTask
	for rows.Next() {
		var task model.GetTask
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			http.Error(w, "Ошибка при чтении результатов", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []model.GetTask{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]model.GetTask{"tasks": tasks})
}

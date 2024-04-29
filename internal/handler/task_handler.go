package handler

import (
	"database/sql"
	"net/http"
)

func TaskHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		AddTask(db, w, r)

	case "GET":
		GetTaskByID(db, w, r)

	case "PUT":
		UpdateTask(db, w, r)

	case "DELETE":
		// Здесь будет логика для удаления задачи
	default:
		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
	}
}

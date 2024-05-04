package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

// DeleteTask удаление задачи
func DeleteTask(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, `{"error":"ID должен быть числом"}`, http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка при удалении задачи: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

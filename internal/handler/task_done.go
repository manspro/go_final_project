package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// MarkTaskDone - отметка о выполнении задачи
func MarkTaskDone(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	var repeat, date string
	err := db.QueryRow("SELECT repeat, date FROM scheduler WHERE id = ?", id).Scan(&repeat, &date)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка базы данных: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if repeat == "" { // удаляем одноразовую задачу
		_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при удалении задачи: %v"}`, err), http.StatusInternalServerError)
			return
		}
	} else { // периодическая задача, обновляем дату
		nextDate, err := NextDate(time.Now(), date, repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при расчёте следующей даты: %v"}`, err), http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при обновлении задачи: %v"}`, err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

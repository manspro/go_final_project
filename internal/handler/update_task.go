package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/internal/model"
	"go_final_project/internal/utils"
	"net/http"
	"strconv"
	"time"
)

func UpdateTask(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var task model.GetTask
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error":"Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	if task.ID == "" || task.Title == "" || task.Date == "" {
		http.Error(w, `{"error":"Необходимые поля не заполнены"}`, http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(task.ID); err != nil {
		http.Error(w, `{"error":"ID должен быть числом"}`, http.StatusBadRequest)
		return
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM scheduler WHERE id = ?", task.ID).Scan(&count)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка базы данных: %v"}`, err), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
		return
	}

	if task.Repeat != "" && !utils.IsValidRepeatFormat(task.Repeat) {
		http.Error(w, `{"error":"Неверный формат повторения"}`, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Ошибка при обновлении задачи: %v"}`, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

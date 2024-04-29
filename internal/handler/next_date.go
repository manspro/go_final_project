package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату на основе правила повторения
func NextDate(now time.Time, dateString string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения пусто")
	}

	date, err := time.Parse("20060102", dateString)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %s", dateString)
	}

	// Проверяем, прошла ли уже исходная дата
	if date.Before(now) {
		// Если дата прошла, начинаем считать от сегодняшнего дня
		date = now
	}

	switch {
	case strings.HasPrefix(repeat, "d "):
		days, err := strconv.Atoi(strings.TrimSpace(repeat[2:]))
		if err != nil || days < 1 {
			return "", fmt.Errorf("неверный интервал дней")
		}
		// Добавляем дни, пока не получим дату в будущем
		nextDate := date.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат повтора")
	}
}

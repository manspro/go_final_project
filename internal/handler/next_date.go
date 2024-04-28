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

	switch {
	case strings.HasPrefix(repeat, "d "):
		days, err := strconv.Atoi(strings.TrimSpace(repeat[2:]))
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("неверный интервал дней")
		}
		nextDate := date.AddDate(0, 0, days)
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil

	case repeat == "y":
		nextDate := date.AddDate(1, 0, 0)
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		// Обработка високосного года если исходная дата 29 февраля
		if date.Month() == time.February && date.Day() == 29 {
			if nextDate.Month() != time.February || nextDate.Day() != 29 {
				nextDate = time.Date(nextDate.Year(), time.March, 1, 0, 0, 0, 0, nextDate.Location())
			}
		}
		return nextDate.Format("20060102"), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат повтора")
	}
}

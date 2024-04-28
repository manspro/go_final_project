package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func SetupDB() *sql.DB {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	createDB(db)

	return db
}

func getDBPath() string {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	return dbFile
}

func createDB(db *sql.DB) {
	createTabSQL := `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    title TEXT NOT NULL,
    comment TEXT,
    repeat TEXT);`

	_, err := db.Exec(createTabSQL)
	if err != nil {
		log.Fatalf("Не удалось создать таблицу: %v", err)
	}

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date)`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		log.Fatalf("Не удалось создать индекс: %v", err)
	}
}

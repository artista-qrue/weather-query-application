package storage

import (
	"database/sql"
	"log"
	"time"
)

type QueryStorage struct {
	db *sql.DB
}

func NewQueryStorage(db *sql.DB) *QueryStorage {
	return &QueryStorage{db: db}
}

func InitDB(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS weather_queries (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        location TEXT,
        service_1_temperature REAL,
        service_2_temperature REAL,
        request_count INTEGER,
        created_at DATETIME
    );`
	_, err := db.Exec(query)
	return err
}

func (s *QueryStorage) SaveQuery(location string, temp1, temp2 float64, requestCount int) {
	query := `
    INSERT INTO weather_queries (location, service_1_temperature, service_2_temperature, request_count, created_at)
    VALUES (?, ?, ?, ?, ?);`
	_, err := s.db.Exec(query, location, temp1, temp2, requestCount, time.Now())
	if err != nil {
		log.Printf("failed to save query: %v", err)
	}
}

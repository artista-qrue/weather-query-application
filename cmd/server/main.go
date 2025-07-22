package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"weather-query-application/internal/aggregator"
	"weather-query-application/internal/config"
	"weather-query-application/internal/handler"
	"weather-query-application/internal/service"
	"weather-query-application/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("sqlite3", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed to close database: %v", err)
		}
	}(db)

	if err := storage.InitDB(db); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	weatherService := service.NewWeatherService(cfg.WeatherAPI.WeatherAPIKey, cfg.WeatherAPI.WeatherStackKey)
	queryStorage := storage.NewQueryStorage(db)
	requestAggregator := aggregator.NewRequestAggregator(weatherService, queryStorage)

	weatherHandler := handler.NewWeatherHandler(requestAggregator)

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", weatherHandler.GetWeather)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("server starting on port %d", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-query-application/internal/aggregator"
	"weather-query-application/internal/service"
	"weather-query-application/internal/storage"
)

func TestWeatherHandler(t *testing.T) {
	weatherService := service.NewWeatherService("dummy_key", "dummy_key")
	queryStorage := storage.NewQueryStorage(nil)
	requestAggregator := aggregator.NewRequestAggregator(weatherService, queryStorage)
	handler := NewWeatherHandler(requestAggregator)

	t.Run("returns error when location is missing", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/weather", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		h := http.HandlerFunc(handler.GetWeather)
		h.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

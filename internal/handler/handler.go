package handler

import (
	"encoding/json"
	"net/http"
	"weather-query-application/internal/aggregator"
)

type WeatherHandler struct {
	aggregator *aggregator.RequestAggregator
}

func NewWeatherHandler(aggregator *aggregator.RequestAggregator) *WeatherHandler {
	return &WeatherHandler{aggregator: aggregator}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("q")
	if location == "" {
		http.Error(w, "location is required", http.StatusBadRequest)
		return
	}

	result, err := h.aggregator.GetWeather(location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

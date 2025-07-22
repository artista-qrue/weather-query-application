package aggregator

import (
	"testing"
	"time"
	"weather-query-application/internal/service"
	"weather-query-application/internal/storage"
)

func TestRequestAggregator(t *testing.T) {
	weatherService := service.NewWeatherService("dummy_key", "dummy_key")
	queryStorage := storage.NewQueryStorage(nil)
	aggregator := NewRequestAggregator(weatherService, queryStorage)

	t.Run("groups requests for the same location", func(t *testing.T) {
		location := "Istanbul"
		go func() {
			_, err := aggregator.GetWeather(location)
			if err != nil {
				t.Errorf("expected 1 request in group, got %d", aggregator.groups[location].requestCount)
			}
		}()
		go func() {
			_, err := aggregator.GetWeather(location)
			if err != nil {
				t.Errorf("expected 1 request in group, got %d", aggregator.groups[location].requestCount)
			}
		}()

		time.Sleep(1 * time.Second)
		aggregator.mutex.Lock()
		if aggregator.groups[location].requestCount != 1 {
			t.Errorf("expected 1 request in group, got %d", aggregator.groups[location].requestCount)
		}
		aggregator.mutex.Unlock()
	})
}

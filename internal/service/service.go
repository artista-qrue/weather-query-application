package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"weather-query-application/model"
)

type WeatherService struct {
	weatherAPIKey   string
	weatherStackKey string
	client          *http.Client
}

func NewWeatherService(weatherAPIKey, weatherStackKey string) *WeatherService {
	return &WeatherService{
		weatherAPIKey:   weatherAPIKey,
		weatherStackKey: weatherStackKey,
		client:          &http.Client{},
	}
}

func (s *WeatherService) GetAverageTemperature(location string) (float64, error) {
	var wg sync.WaitGroup
	temp1Chan := make(chan float64, 1)
	temp2Chan := make(chan float64, 1)
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		temp, err := s.getWeatherAPITemperature(location)
		if err != nil {
			errChan <- err
			return
		}
		temp1Chan <- temp
	}()

	go func() {
		defer wg.Done()
		temp, err := s.getWeatherStackTemperature(location)
		if err != nil {
			errChan <- err
			return
		}
		temp2Chan <- temp
	}()

	wg.Wait()
	close(errChan)
	close(temp1Chan)
	close(temp2Chan)

	if len(errChan) > 0 {
		return 0, <-errChan
	}

	temp1 := <-temp1Chan
	temp2 := <-temp2Chan

	return (temp1 + temp2) / 2, nil
}

func (s *WeatherService) getWeatherAPITemperature(location string) (float64, error) {
	//todo it can be managedc via env. empire bless you
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=1&aqi=no&alerts=no", s.weatherAPIKey, location)
	resp, err := s.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data model.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}

func (s *WeatherService) getWeatherStackTemperature(location string) (float64, error) {
	url := fmt.Sprintf("http://api.weatherstack.com/current?access_key=%s&query=%s", s.weatherStackKey, location)
	resp, err := s.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data model.WeatherStackResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return float64(data.Current.Temperature), nil
}

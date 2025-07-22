package aggregator

import (
	"sync"
	"time"
	"weather-query-application/internal/service"
	"weather-query-application/internal/storage"
	"weather-query-application/model"
)

type requestGroup struct {
	waitGroup    sync.WaitGroup
	result       *model.WeatherResponse
	err          error
	requestCount int
}

type RequestAggregator struct {
	weatherService *service.WeatherService
	queryStorage   *storage.QueryStorage
	groups         map[string]*requestGroup
	mutex          sync.Mutex
}

func NewRequestAggregator(weatherService *service.WeatherService, queryStorage *storage.QueryStorage) *RequestAggregator {
	return &RequestAggregator{
		weatherService: weatherService,
		queryStorage:   queryStorage,
		groups:         make(map[string]*requestGroup),
	}
}

func (a *RequestAggregator) GetWeather(location string) (*model.WeatherResponse, error) {
	a.mutex.Lock()
	group, ok := a.groups[location]
	if !ok {
		group = &requestGroup{}
		group.waitGroup.Add(1)
		a.groups[location] = group
		a.mutex.Unlock()

		go a.processGroup(location, group)
	} else {
		group.requestCount++
		a.mutex.Unlock()
	}

	group.waitGroup.Wait()
	return group.result, group.err
}

func (a *RequestAggregator) processGroup(location string, group *requestGroup) {
	defer func() {
		a.mutex.Lock()
		delete(a.groups, location)
		a.mutex.Unlock()
		group.waitGroup.Done()
	}()

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			a.fetchAndRespond(location, group)
			return
		default:
			if group.requestCount >= 10 {
				a.fetchAndRespond(location, group)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (a *RequestAggregator) fetchAndRespond(location string, group *requestGroup) {
	temp1, err := a.weatherService.GetAverageTemperature(location)
	if err != nil {
		group.err = err
		return
	}

	temp2, err := a.weatherService.GetAverageTemperature(location)
	if err != nil {
		group.err = err
		return
	}

	avgTemp := (temp1 + temp2) / 2
	group.result = &model.WeatherResponse{
		Location:    location,
		Temperature: avgTemp,
	}

	go a.queryStorage.SaveQuery(location, temp1, temp2, group.requestCount+1)
}

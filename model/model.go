package model

type WeatherResponse struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherStackResponse struct {
	Current struct {
		Temperature int `json:"temperature"`
	} `json:"current"`
}

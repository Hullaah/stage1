package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io"
	"os"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := fn(w, r); err != nil {
        http.Error(w, err.Error(), 500)
    }
}

type Location  struct {
	City string
	Latitude float32
	Longitude float32
}

func getClientLocation(clientIp string) (*Location, error) {
	ipApiKey := os.Getenv("IP_API_KEY")
	resp, err := http.Get(fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s&ip=%s", ipApiKey, clientIp))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	location := new(Location)
	json.Unmarshal(body, location)
	return location, nil
}

func getLocationWeather(location *Location) (float32, error) {
	weatherApiKey := os.Getenv("WEATHER_API_KEY")
	args := []any{location.Latitude, location.Longitude, weatherApiKey}
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", args...)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var weather struct {Main struct {Temp float32}}
	json.Unmarshal(body, &weather)
	return weather.Main.Temp, nil
}


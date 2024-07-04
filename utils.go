package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

type Location struct {
	City      string
	Latitude  float32
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
	var weather struct{ Main struct{ Temp float32 } }
	json.Unmarshal(body, &weather)
	return weather.Main.Temp, nil
}

func getClientIP(r *http.Request) string {
	// Check the X-Forwarded-For header first
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can contain multiple IPs, the first one is the client IP
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// If X-Forwarded-For is not set, check the X-Real-IP header
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fallback to the remote address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

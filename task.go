package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"io"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		visitorName := r.URL.Query().Get("visitor_name")
		addr, _ := net.LookupHost(r.URL.Hostname())
		ipApiKey := os.Getenv("IP_API_KEY")
		weatherApiKey := os.Getenv("WEATHER_API_KEY")
		clientIp := addr[0]
		w.Header().Set("Content-Type", "application/json")

		geolocationResp, _ := http.Get(fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s&ip=%s", ipApiKey, clientIp))
		defer geolocationResp.Body.Close()
		body, _ := io.ReadAll(geolocationResp.Body)
		var location struct {City string; Latitude float32; Longitude float32}
		json.Unmarshal(body, &location)

		weatherResp, _ := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", location.Latitude, location.Longitude, weatherApiKey))
		defer weatherResp.Body.Close()
		body, _ = io.ReadAll(weatherResp.Body)
		var weather struct {Main struct {Temp float32}}
		json.Unmarshal(body, &weather)

		greeting := fmt.Sprintf("Hello, %s!, the temperature is %.2f degrees Celcius in %s", visitorName, weather.Main.Temp, location.City)
		out, _ := json.Marshal(struct {
			ClientIp string `json:"client_ip"`
			Location string `json:"location"`
			Greeting string `json:"greeting"`
		}{addr[0], location.City, greeting})
		w.Write(out)
	})

	http.ListenAndServe(":8080", nil)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	http.Handle("/api/hello", appHandler(view))
	http.ListenAndServe(":8080", nil)
}

func view(w http.ResponseWriter, r *http.Request) error {
	clientIp := getClientIP(r)

	location, err := getClientLocation(clientIp)
	if err != nil {
		return err
	}

	temperature, err := getLocationWeather(location)
	if err != nil {
		return err
	}

	visitorName := r.URL.Query().Get("visitor_name")
	greeting := fmt.Sprintf("Hello, %s!, the temperature is %.2f degrees Celcius in %s", visitorName, temperature, location.City)
	out, _ := json.Marshal(struct {
		ClientIp string `json:"client_ip"`
		Location string `json:"location"`
		Greeting string `json:"greeting"`
	}{clientIp, location.City, greeting})

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
	return nil
}

package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Response struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

var (
	ErrInvalidZipCode  = errors.New("invalid zipcode")
	ErrAddressNotFound = errors.New("can not find zipcode")
)

func main() {
	weatherApiKey := os.Getenv("WEATHER_API_KEY")
	viaCepClient := NewViaCepClient("https://viacep.com.br/ws/")
	weatherClient := NewWeatherClient("https://api.weatherapi.com/v1/current.json", weatherApiKey)

	server := makeServer(viaCepClient, weatherClient)
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}

func makeServer(viaCepClient *ViaCepClient, weatherClient *WeatherClient) *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("GET /weather", func(writer http.ResponseWriter, request *http.Request) {
		postalCode := request.URL.Query().Get("postalCode")

		var regex = regexp.MustCompile(`\D`)
		postalCode = regex.ReplaceAllString(postalCode, "")
		if len(postalCode) != 8 {
			responseError(writer, http.StatusUnprocessableEntity, ErrInvalidZipCode)
			return
		}

		viaCepResponse, err := viaCepClient.GetAddress(postalCode)
		if err != nil {
			responseError(writer, http.StatusInternalServerError, err)
			return
		}

		if viaCepResponse.Error != nil {
			responseError(writer, http.StatusNotFound, ErrAddressNotFound)
			return
		}

		weatherResponse, err := weatherClient.GetWeather(viaCepResponse.Localidade)
		if err != nil {
			responseError(writer, http.StatusInternalServerError, err)
			return
		}

		celsius := weatherResponse.Current.TempC
		response := Response{
			Celsius:    round(celsius),
			Fahrenheit: round(celsius*1.8 + 32),
			Kelvin:     round(celsius + 273),
		}
		_ = json.NewEncoder(writer).Encode(response)
	})
	return server
}

func responseError(writer http.ResponseWriter, status int, err error) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(ErrorResponse{Error: err.Error()})
}

func round(value float64) float64 {
	intValue := int(value * 100)
	return float64(intValue) / 100
}

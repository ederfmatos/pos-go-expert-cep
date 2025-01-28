package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_GetWeather(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{postalCode}/json", func(writer http.ResponseWriter, request *http.Request) {
		postalCode := request.PathValue("postalCode")
		encoder := json.NewEncoder(writer)
		if postalCode == "00000000" {
			errMessage := "not found"
			_ = encoder.Encode(ViaCepResponse{Error: &errMessage})
			return
		} else {
			_ = encoder.Encode(ViaCepResponse{Localidade: "City"})
		}
		writer.WriteHeader(200)
	})
	mux.HandleFunc("GET /weather", func(writer http.ResponseWriter, request *http.Request) {
		_ = json.NewEncoder(writer).Encode(WeatherResponse{Current: WeatherCurrent{TempC: 30}})
		writer.WriteHeader(200)
	})

	externalService := httptest.NewServer(mux)
	defer externalService.Close()

	server := makeServer(NewViaCepClient(externalService.URL), NewWeatherClient(externalService.URL+"/weather", ""))
	httpServer := httptest.NewServer(server)

	t.Run("should return 422 if postal code is not provided", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/weather", nil)
		recorder := httptest.NewRecorder()

		server.ServeHTTP(recorder, request)
		requireEqual(t, http.StatusUnprocessableEntity, recorder.Code, "StatusCode")
		requireContains(t, recorder.Body.String(), ErrInvalidZipCode.Error())
	})

	t.Run("should return 422 if invalid postal code is not provided", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/weather?postalCode=1234567890", nil)
		recorder := httptest.NewRecorder()

		server.ServeHTTP(recorder, request)
		requireEqual(t, http.StatusUnprocessableEntity, recorder.Code, "StatusCode")
		requireContains(t, recorder.Body.String(), ErrInvalidZipCode.Error())
	})

	t.Run("should return 404 if postal code is not valid for a city", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/weather?postalCode=00000000", nil)
		recorder := httptest.NewRecorder()

		server.ServeHTTP(recorder, request)
		requireEqual(t, http.StatusNotFound, recorder.Code, "StatusCode")
		requireContains(t, recorder.Body.String(), ErrAddressNotFound.Error())
	})

	t.Run("should return 200 on success", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/weather?postalCode=12345678", nil)
		recorder := httptest.NewRecorder()

		server.ServeHTTP(recorder, request)
		requireEqual(t, http.StatusOK, recorder.Code, "StatusCode")

		var response Response
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		requireNoError(t, err, "Unmarshal response")

		requireEqual(t, 30.0, response.Celsius, "Celsius")
		requireEqual(t, 303.0, response.Kelvin, "Kelvin")
		requireEqual(t, 86.0, response.Fahrenheit, "Fahrenheit")
	})

	defer httpServer.Close()
}

func requireEqual(t *testing.T, expected, actual interface{}, message string) {
	if expected != actual {
		t.Fatalf("expected %v but got %v - [%s]", expected, actual, message)
	}
}

func requireContains(t *testing.T, value, contains string) {
	if !strings.Contains(value, contains) {
		t.Fatalf("expected %s contains %s", value, contains)
	}
}

func requireNoError(t *testing.T, err error, message string) {
	if err != nil {
		t.Fatalf("expected no error %s - [%s]", err, message)
	}
}

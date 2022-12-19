package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	apiKey = "903d17b3c177baff05558f30dc1da601"
	apiURL = "https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=" + apiKey
)

// City represents a city with its current weather and temperature.
type City struct {
	Name        string  `json:"name"`
	Temperature float64 `json:"temperature"`
	Conditions  string  `json:"conditions"`
}

// WeatherData represents the response from the OpenWeatherMap API.
type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Pressure float64 `json:"pressure"`
		Humidity float64 `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
}

// WeatherAPI represents the weather API.
type WeatherAPI struct {
	cache *cache.Cache
}

// NewWeatherAPI creates a new instance of the weather API.
func NewWeatherAPI() *WeatherAPI {
	return &WeatherAPI{
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

// GetCitiesWeather retrieves the weather data for the given list of cities.
func (api *WeatherAPI) GetCitiesWeather(cities []string) ([]City, error) {
	var data []City

	for _, city := range cities {
		// Check the cache to see if the weather data for this city is already available.
		if item, found := api.cache.Get(city); found {
			// If the data is in the cache, unmarshal it into a City struct.
			var c City
			if err := json.Unmarshal(item.([]byte), &c); err != nil {
				return nil, err
			}
			data = append(data, c)
			continue
		}

		// If the data is not in the cache, make a request to the weather API.
		resp, err := http.Get(fmt.Sprintf(apiURL, city))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Decode the response from the weather API.
		var weatherData WeatherData
		if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
			return nil, err
		}

		// Create a new City struct with the weather data.
		c := City{
			Name:        weatherData.Name,
			Temperature: weatherData.Main.Temp,
			Conditions:  weatherData.Weather[0].Description,
		} // Marshal the City struct into JSON and add it to the cache.
		item, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}
		api.cache.Set(city, item, cache.DefaultExpiration)

		data = append(data, c)
	}

	return data, nil
}

func main() {
	api := NewWeatherAPI()
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		// Set the content type to JSON.
		w.Header().Set("Content-Type", "application/json")

		// Get the list of cities from the query parameters.
		cities := r.URL.Query()["city"]
		if len(cities) == 0 {
			http.Error(w, "missing city query parameter", http.StatusBadRequest)
			return
		}

		// Retrieve the weather data for the cities.
		data, err := api.GetCitiesWeather(cities)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the weather data as a JSON response.
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}

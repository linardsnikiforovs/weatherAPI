#WeatherAPI
WeatherAPI is a Go package that provides a simple interface for retrieving weather data for a list of cities. It uses the OpenWeatherMap API to retrieve the data and has a built-in cache to store the results for a given period of time.

##Installation
To install WeatherAPI, use go get:

```Code :-go get github.com/[YOUR_USERNAME]/weatherapi```

##Usage
```import "github.com/[YOUR_USERNAME]/weatherapi"

api := weatherapi.NewWeatherAPI()

cities, err := api.GetCitiesWeather([]string{"Paris", "London", "New York"})
if err != nil {
	// Handle error
}

// cities will contain a list of City structs with the weather data for each city
```

##API
###func NewWeatherAPI() *WeatherAPI
NewWeatherAPI creates a new instance of the WeatherAPI.

###func (api *WeatherAPI) GetCitiesWeather(cities []string) ([]City, error)
GetCitiesWeather retrieves the weather data for the given list of cities. It returns a slice of City structs containing the name, temperature, and weather conditions for each city. If there is an error retrieving the data, it returns an error.

###Structs
###type City
City represents a city with its current weather and temperature.

```type City struct {
	Name        string  // Name of the city
	Temperature float64 // Current temperature in Celsius
	Conditions  string  // Description of the current weather conditions
}
```

##Example
Here is an example of how you might use the WeatherAPI to retrieve the weather data for a list of cities and serve it as a JSON API:

```package main

import (
	"encoding/json"
	"net/http"

	"github.com/[YOUR_USERNAME]/weatherapi"
)

func main() {
	api := weatherapi.NewWeatherAPI()
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		// Set the content type to JSON.
		w.Header().Set("Content-Type", "application/json")

		// Get the list of cities from the query parameters.
		cities := r.URL.Query()["city"]
		if len(cities) == 0 {
			http.Error(w, "No cities provided", http.StatusBadRequest)
			return
		}

		// Retrieve the weather data for the cities.
		data, err := api.GetCitiesWeather(cities)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Marshal the data into JSON and write it to the response.
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
```

##Configuration
The WeatherAPI uses the OpenWeatherMap API to retrieve the weather data. You will need to provide your own API key by setting the apiKey constant in the package. You can obtain an API key by creating a free account on the OpenWeatherMap website.

The WeatherAPI also has a built-in cache to store the results for a given period of time. The cache expiration time can be configured by modifying the parameters passed to the cache.New function in the NewWeatherAPI function.

##Limitations
The OpenWeatherMap API has some limitations on the number of requests that can be made per minute. Be sure to read the API documentation to understand the limitations and how to avoid going over the rate limit.

##License
WeatherAPI is released under the MIT License. See the LICENSE file for more details.

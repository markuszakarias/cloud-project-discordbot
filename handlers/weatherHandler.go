package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"projectGroup23/database"
	"projectGroup23/structs"
	"projectGroup23/utils"
)

// struct used to retrieved data from api
var weatherForecast structs.WeatherForecasts

// const for cache duration
const weatherForecastDur = 100

// getWeatherForecastAndIP
func getWeatherForecastAndIP(location string) (structs.WeatherForecasts, error) {
	var err error = nil
	fmt.Println("API call made!") // for debugging

	apikey := os.Getenv("WEATHER_KEY")
	wf, err := http.Get("https://api.openweathermap.org/data/2.5/forecast/daily?q=" + location + "&units=metric&cnt=1&appid=" + apikey)

	if err != nil {
		return weatherForecast, err
	}
	output, err := ioutil.ReadAll(wf.Body)
	if err != nil {
		return weatherForecast, err
	}
	jsonRes := string(output)

	weatherForecast = utils.PopulateWeatherForecast(jsonRes, 1)

	// return the populated object
	// cache the data retrieved from API
	err = storeWeatherForecastAndIP(weatherForecast, location)

	// return the populated object
	return weatherForecast, err
}

// WeatherForecastMainHandler
func WeatherForecastMainHandler(location string) (structs.WeatherForecasts, error) {
	var err error = nil
	// use function to retrieve cached newsletter
	wf := getStoredWeatherForecast()

	// check if the interface is null
	if wf.Forecasts == nil || wf.Forecasts[0].City != location {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		database.DeleteObjectFromFirestore(database.StoredWeatherForecast.FirestoreID)
		wf, err = getWeatherForecastAndIP(location)
	}

	return wf, err
}

// storeWeatherForecastAndIP
func storeWeatherForecastAndIP(resp structs.WeatherForecasts, ipLoc string) error {
	// populate struct with data to be cached
	database.StoredWeatherForecast = structs.StoredWeatherForecast{
		WeatherForecasts: resp,
		IPLocation:       ipLoc,
		StoreTime:        time.Now(),
		StoreRefresh:     weatherForecastDur,
	}
	// save the object on firestore
	err := database.SaveWeatherForecastToFirestore(&database.StoredWeatherForecast)
	return err
}

// getStoredWeatherForecast
func getStoredWeatherForecast() structs.WeatherForecasts {
	if database.StoredWeatherForecast.WeatherForecasts.Forecasts == nil {
		return structs.WeatherForecasts{}
	}
	database.StoredWeatherForecast.StoreRefresh -= time.Since(database.StoredWeatherForecast.StoreTime).Seconds()
	database.StoredWeatherForecast.StoreTime = time.Now()
	database.UpdateTimeFirestore(database.StoredWeatherForecast.FirestoreID, database.StoredWeatherForecast.StoreTime, database.StoredWeatherForecast.StoreRefresh)
	fmt.Println(database.StoredWeatherForecast.StoreRefresh)
	if database.StoredWeatherForecast.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(database.StoredWeatherForecast.FirestoreID)
		return structs.WeatherForecasts{}
	}
	return database.StoredWeatherForecast.WeatherForecasts
}

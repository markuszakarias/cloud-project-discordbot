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

	fmt.Print("Handler was called!")

	storedwf, err := database.CheckWeatherForecastsOnFirestore(location)
	if err != nil {
		fmt.Println(err)
	}

	wf := getStoredWeatherForecast(storedwf)

	// check if the interface is null
	if wf.Forecasts == nil || wf.Forecasts[0].City != location {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
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
func getStoredWeatherForecast(storedwf structs.StoredWeatherForecast) structs.WeatherForecasts {
	if storedwf.WeatherForecasts.Forecasts == nil {
		return storedwf.WeatherForecasts
	}
	storedwf.StoreRefresh -= time.Since(storedwf.StoreTime).Seconds()
	storedwf.StoreTime = time.Now()
	database.UpdateTimeFirestore(storedwf.FirestoreID, storedwf.StoreTime, storedwf.StoreRefresh)
	fmt.Println(storedwf.StoreRefresh)
	if storedwf.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(storedwf.FirestoreID)
		return structs.WeatherForecasts{}
	}
	return storedwf.WeatherForecasts
}

package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"projectGroup23/database"
	"projectGroup23/structs"
	"projectGroup23/utils"

	"github.com/tidwall/gjson"
)

// Struct used to handle data in database management system
var weatherForecast structs.WeatherForecasts

// Const for database storage duration
const weatherForecastDur = 100

// GetWeatherForecastAndIP - Requests weather forecast and ip location from two api's
// this call is only done when no stored data exists at startup
// and when a stored object is deleted after timeout
func GetWeatherForecastAndIP(location string) (structs.WeatherForecasts, error) {
	// Get api key from env variable
	apikey := os.Getenv("WEATHER_KEY")

	resp, err := http.Get("https://api.openweathermap.org/data/2.5/forecast/daily?q=" + location + "&units=metric&cnt=1&appid=" + apikey)
	if err != nil {
		return weatherForecast, err
	}

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return weatherForecast, err
	}

	jsonRes := string(output)

	cod := gjson.Get(jsonRes, "cod").String()

	if cod == "404" { // if input city does not exist
		return weatherForecast, errors.New("city '" + location + "' not found")
	}
	if err != nil {
		return weatherForecast, err
	}

	// Populates object with JSON response
	weatherForecast = utils.PopulateWeatherForecast(jsonRes, 1)

	// Store the data retrieved from API
	err = storeWeatherForecastAndIP(weatherForecast, location)

	return weatherForecast, err
}

// WeatherForecastMainHandler - Main handler for the !weather command
func WeatherForecastMainHandler(location string) (structs.WeatherForecasts, error) {
	var err error

	// Checks if the requested data exists in database
	stored, err := database.CheckWeatherForecastsOnFirestore(location)
	if err != nil {
		fmt.Println(err)
	}

	// Retrieving possible stored data
	wf := getStoredWeatherForecast(stored)

	// Checks if it exists stored data
	if wf.Forecasts == nil || wf.Forecasts[0].City != location {
		wf, err = GetWeatherForecastAndIP(location)
	}

	return wf, err
}

// storeWeatherForecastAndIP - Stores a WeatherForecasts and IPLocation object in the database
func storeWeatherForecastAndIP(resp structs.WeatherForecasts, ipLoc string) error {
	// Populate struct with data to be stored
	database.StoredWeatherForecast = structs.StoredWeatherForecast{
		WeatherForecasts: resp,
		IPLocation:       ipLoc,
		StoreTime:        time.Now(),
		StoreRefresh:     weatherForecastDur,
	}
	// Store the object
	err := database.SaveWeatherForecastToFirestore(&database.StoredWeatherForecast)
	return err
}

// getStoredWeatherForecast - Updates timestamps in database storage and retrieves matching object to request
func getStoredWeatherForecast(stored structs.StoredWeatherForecast) structs.WeatherForecasts {
	// Checks if it exists a stored response
	if stored.WeatherForecasts.Forecasts == nil {
		return stored.WeatherForecasts
	}

	// Calculates timestamps and duration stored in database
	stored.StoreRefresh -= time.Since(stored.StoreTime).Seconds()
	stored.StoreTime = time.Now()

	// Updates new timestamp and duration to Firestore object
	database.UpdateTimeFirestore(stored.FirestoreID, stored.StoreTime, stored.StoreRefresh)

	// If the object storage timer is timed out the object is deleted and then renewed when the next command is called
	if stored.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(stored.FirestoreID)
		return structs.WeatherForecasts{}
	}
	return stored.WeatherForecasts
}

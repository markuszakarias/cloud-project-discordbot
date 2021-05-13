package handlers

import (
	"encoding/json"
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

// struct used to retrieved IP location from api
var ipAddress structs.IPLocation

// const for cache duration
const weatherForecastDur = 100

func getIPLocation() (structs.IPLocation, error) {
	resp, err := http.Get("https://ipwhois.app/json/")
	if err != nil {
		return ipAddress, err
		//fmt.Errorf("Error in response: ", err.Error())
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ipAddress)
	if err != nil {
		return ipAddress, err
		//fmt.Errorf("Error in JSON decoding: ", err.Error())
	}

	return ipAddress, nil
}

// getNewsletters - gets all the newsletters from the api
// this call is only done when no cached data exists at startup
// and when a cached object is deleted after timeout
// TODO - security on API key
// TODO - better error handling
func getWeatherForecastAndIP() (structs.WeatherForecasts, error) {
	var err error = nil
	fmt.Println("API call made!") // for debugging

	ipAddress, err = getIPLocation()
	if err != nil {
		return weatherForecast, err
	}

	apikey := os.Getenv("WEATHER_KEY")
	wf, err := http.Get("https://api.openweathermap.org/data/2.5/forecast/daily?q=Oslo&units=metric&cnt=1&appid=" + os.Getenv(apikey))

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
	err = storeWeatherForecastAndIP(weatherForecast, ipAddress)

	// return the populated object
	return weatherForecast, err
}

// TestEndpoint - just for development, testing that the functionality works correctly
// TODO - remove when not needed anymore
func WeatherForecastMainHandler() (structs.WeatherForecasts, error) {
	var err error = nil
	// use function to retrieve cached newsletter
	wf := getStoredWeatherForecast()

	// check if the interface is null
	if wf.Forecasts == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		wf, err = getWeatherForecastAndIP()
	}

	return wf, err
}

//TODO Look into merging storeWeatherForecastAndIP and saveWeatherForecastToFirestore

// cacheNewsLetter - caches a NewsLetter object to a cache object
// will be stored in firestore
func storeWeatherForecastAndIP(resp structs.WeatherForecasts, ipLoc structs.IPLocation) error {
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

// getCachedNewsLetter - used on endpoint to retrieve the cached newsletter
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
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

package caching

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"projectGroup23/database"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// struct used for the cached data
var c_weather CachedWeatherForecast

// struct used to retrieved data from api
var weather interface{}

// struct used to retrieved IP location from api
var ipAdress IPLocation

// const for cache duration
const c_weatherforecast_dur = 100

type CachedWeatherForecast struct {
	WeatherForecasts interface{}
	IPLocation       IPLocation
	CachedTime       time.Time
	CachedRefresh    float64
	firestoreID      string
}

type IPLocation struct {
	City string `json:"city"`
}

func getIPLocation() IPLocation {
	resp, err := http.Get("https://ipwhois.app/json/")
	if err != nil {
		fmt.Errorf("Error in response: ", err.Error())
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ipAdress)
	if err != nil {
		fmt.Errorf("Error in JSON decoding: ", err.Error())
	}

	return ipAdress
}

// getNewsletters - gets all the newsletters from the api
// this call is only done when no cached data exists at startup
// and when a cached object is deleted after timeout
// TODO - security on API key
// TODO - better error handling
func getWeatherForecastAndIP() interface{} {
	fmt.Println("API call made!") // for debugging

	ipAdress = getIPLocation()

	wf, err := http.Get("https://api.openweathermap.org/data/2.5/forecast/daily?q=Oslo&units=metric&cnt=1&appid=f6a8e67b1a5f1d5be2bffe4d461cc155")

	if err != nil {
		fmt.Println(err)
	}
	err = json.NewDecoder(wf.Body).Decode(&weather)
	if err != nil {
		fmt.Println(err)
	}
	// cache the data retrieved from API
	cacheWeatherForecastAndIP(weather, ipAdress)

	// return the populated object
	return weather
}

// TestEndpoint - just for development, testing that the functionality works correctly
// TODO - remove when not needed anymore
func WeatherForecastTest(w http.ResponseWriter, r *http.Request) {
	// use function to retrieve cached newsletter
	wf := getCachedWeatherForecast()

	// check if the interface is null
	if wf == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		wf = getWeatherForecastAndIP()
	}
	err := json.NewEncoder(w).Encode(wf)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// cacheNewsLetter - caches a NewsLetter object to a cache object
// will be stored in firestore
func cacheWeatherForecastAndIP(resp interface{}, iploc IPLocation) {
	fmt.Println("caching!")
	// populate struct with data to be cached
	c_weather = CachedWeatherForecast{
		WeatherForecasts: resp,
		IPLocation:       iploc,
		CachedTime:       time.Now(),
		CachedRefresh:    c_weatherforecast_dur,
	}
	// save the object on firestore
	saveWeatherForecastToFirestore(&c_weather)
}

// saveNewsLetterToFirestore - saves an object to firestore
func saveWeatherForecastToFirestore(c_save *CachedWeatherForecast) {
	fmt.Println("Savetofirestore!")
	doc, _, err := database.Client.Collection("cached_resp").Add(database.Ctx, *c_save)
	c_save.firestoreID = doc.ID     // storing firestore ID for later use
	fmt.Println(c_save.firestoreID) // confirming the storage of document ID
	if err != nil {
		fmt.Println(err)
	}
}

// GetCachedNewsLetterFromFirestore - global function that runs at startup
// gets all the cached data from firestore
func GetCachedWeatherForecastFromFirestore() {
	iter := database.Client.Collection("cached_resp").Documents(database.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		doc.DataTo(&c_news)
		c_weather.firestoreID = doc.Ref.ID // matching the firestore ID with the one stored
	}
}

// getCachedNewsLetter - used on endpoint to retrieve the cached newsletter
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
func getCachedWeatherForecast() interface{} {
	if c_weather.WeatherForecasts == nil {
		return nil
	}
	c_weather.CachedRefresh -= time.Since(c_weather.CachedTime).Seconds()
	c_weather.CachedTime = time.Now()
	updateCachedTimeOnWeatherForecastFirestore(c_weather.firestoreID, c_weather.CachedTime, c_weather.CachedRefresh)
	fmt.Println(c_weather.CachedRefresh)
	if c_weather.CachedRefresh <= 0 {
		deleteWeatherForecastFromFirestore(c_weather.firestoreID)
		return nil
	}
	return c_weather.WeatherForecasts
}

// deleteNewsLetterFromFirestore - deletes an object in firestore based on firestore ID
func deleteWeatherForecastFromFirestore(firestoreID string) {
	_, err := database.Client.Collection("cached_resp").Doc(firestoreID).Delete(database.Ctx)
	if err != nil {
		fmt.Println(err)
	}
}

// updateCachedTimeOnNewsLetterFirestore - updates the object in firestore
func updateCachedTimeOnWeatherForecastFirestore(firestoreID string, cachedTime time.Time, cachedRefresh float64) {
	_, err := database.Client.Collection("cached_resp").Doc(firestoreID).Update(database.Ctx, []firestore.Update{
		{
			Path:  "CachedTime", // matching specific field in firestore object
			Value: cachedTime,
		},
		{
			Path:  "CachedRefresh", // matching specific field in firestore object
			Value: cachedRefresh,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

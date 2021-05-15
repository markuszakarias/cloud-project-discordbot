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
var newsLetter structs.NewsLetters

// Const for number of newsletters stored in firebase
const newsStored = 3

// Const for database storage duration
const newsLetterDur = 3600

// getNewsletters - Requests all newsletters from the api
// this call is only done when no stored data exists at startup
// and when a stored object is deleted after timeout
func getNewsletters(country string) (structs.NewsLetters, error) {
	// Get api key from env variable
	apikey := os.Getenv("NEWS_KEY")

	resp, err := http.Get("https://newsapi.org/v2/top-headlines?country=" + country + "&apiKey=" + apikey)
	if err != nil {
		return newsLetter, err
	}

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return newsLetter, err
	}

	jsonRes := string(output)

	articles := gjson.Get(jsonRes, "totalResults").Float()

	if articles == 0 { // If input country has no newsletter articles
		return newsLetter, errors.New("country '" + country + "' has no newsletter articles")
	}
	if err != nil {
		return newsLetter, err
	}

	// Populates object with JSON response
	newsLetter = utils.PopulateNewsLetters(newsStored, jsonRes)

	// Store the data retrieved from API
	err = storeNewsLetter(newsLetter, country)

	return newsLetter, err
}

// NewsLetterMainHandler - Main handler for the !newsletter command
func NewsLetterMainHandler(country string) (structs.NewsLetters, error) {
	var err error

	// Checks if the requested data exists in database
	storedNews, err := database.CheckNewsLetterOnFirestore(country)
	if err != nil {
		fmt.Println(err)
	}

	// Retrieving possible stored data
	nws := getStoredNewsLetter(storedNews)

	// Checks if it exists stored data
	if nws.Newsletters == nil || storedNews.Location != country {
		nws, err = getNewsletters(country)
	}

	return nws, err
}

// storeNewsLetter - Stores a NewsLetters object in the database
func storeNewsLetter(resp structs.NewsLetters, country string) error {
	// Populate struct with data to be stored
	database.StoredNewsLetter = structs.StoredNewsLetter{
		NewsLetters:  resp,
		Location:     country,
		StoreTime:    time.Now(),
		StoreRefresh: newsLetterDur,
	}
	// Store the object
	err := database.SaveNewsLetterToFirestore(&database.StoredNewsLetter)
	return err
}

// getStoredNewsLetter - Updates timestamps in database storage and retrieves matching object to request
func getStoredNewsLetter(stored structs.StoredNewsLetter) structs.NewsLetters {
	// Checks if it exists a stored response
	if stored.NewsLetters.Newsletters == nil {
		return structs.NewsLetters{}
	}

	// Calculates timestamps and duration stored in database
	stored.StoreRefresh -= time.Since(stored.StoreTime).Seconds()
	stored.StoreTime = time.Now()

	// Updates new timestamp and duration to Firestore object
	database.UpdateTimeFirestore(stored.FirestoreID, stored.StoreTime, stored.StoreRefresh)

	// If the object storage timer is timed out the object is deleted and then renewed when the next command is called
	if stored.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(stored.FirestoreID)
		return structs.NewsLetters{}
	}
	return stored.NewsLetters
}

package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"projectGroup23/database"
	"projectGroup23/structs"
	"projectGroup23/utils"
)

var newsLetter structs.NewsLetters

// const for cache duration
const newsLetterDur = 100

// getNewsletters - gets all the newsletters from the api
// this call is only done when no cached data exists at startup
// and when a cached object is deleted after timeout
// TODO - security on API key
// TODO - better error handling
func getNewsletters(apikey string) (structs.NewsLetters, error) {
	fmt.Println("API call made!") // for debugging
	resp, err := http.Get("https://newsapi.org/v2/top-headlines?country=no&apiKey=" + apikey)

	if err != nil {
		return newsLetter, err
	}
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return newsLetter, err
	}
	jsonRes := string(output)

	newsLetter = utils.PopulateNewsLetters(3, jsonRes)

	// cache the data retrieved from API
	err = storeNewsLetter(newsLetter)

	return newsLetter, err
}

// TestEndpoint - just for development, testing that the functionality works correctly
// TODO - remove when not needed anymore
func NewsLetterMainHandler(apikey string) (structs.NewsLetters, error) {
	var err error
	fmt.Println("NewsletterTest() was run!")
	// use function to retrieve cached newsletter
	nws := getStoredNewsLetter()

	// check if the interface is null
	if nws.Newsletters == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		nws, err = getNewsletters(apikey)
	}

	return nws, err
}

// storeNewsLetter - stores a NewsLetter object to Firestore
func storeNewsLetter(resp structs.NewsLetters) error {
	// populate struct with data to be store
	database.StoredNewsLetter = structs.StoredNewsLetter{
		NewsLetters:  resp,
		StoreTime:    time.Now(),
		StoreRefresh: newsLetterDur,
	}
	// save the object to firestore
	err := database.SaveNewsLetterToFirestore(&database.StoredNewsLetter)
	return err
}

// getCachedNewsLetter - used on endpoint to retrieve the cached newsletter
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
func getStoredNewsLetter() structs.NewsLetters {
	if database.StoredNewsLetter.NewsLetters.Newsletters == nil {
		return structs.NewsLetters{}
	}
	database.StoredNewsLetter.StoreRefresh -= time.Since(database.StoredNewsLetter.StoreTime).Seconds()
	database.StoredNewsLetter.StoreTime = time.Now()
	database.UpdateTimeFirestore(database.StoredNewsLetter.FirestoreID, database.StoredNewsLetter.StoreTime, database.StoredNewsLetter.StoreRefresh)
	fmt.Println(database.StoredNewsLetter.StoreRefresh)
	if database.StoredNewsLetter.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(database.StoredNewsLetter.FirestoreID)
		return structs.NewsLetters{}
	}
	return database.StoredNewsLetter.NewsLetters
}

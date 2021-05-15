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

var newsLetter structs.NewsLetters

// const for newsletters stored in firebase
const newsStored = 3

// const for cache duration
const newsLetterDur = 100

// getNewsletters - gets all the newsletters from the api
// this call is only done when no cached data exists at startup
// and when a cached object is deleted after timeout
func getNewsletters(country string) (structs.NewsLetters, error) {
	fmt.Println("API call made!") // for debugging

	fmt.Println(country)
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

	if articles == 0 { // if input country has not newsletter articles
		return newsLetter, errors.New("country '" + country + "' has no newsletter articles")
	}

	if err != nil {
		return newsLetter, err
	}

	newsLetter = utils.PopulateNewsLetters(newsStored, jsonRes)

	// cache the data retrieved from API
	err = storeNewsLetter(newsLetter, country)

	return newsLetter, err
}

// NewsLetterMainHandler
func NewsLetterMainHandler(country string) (structs.NewsLetters, error) {
	var err error
	fmt.Println("NewsletterTest() was run!")
	// use function to retrieve cached newsletter

	storedNews, err := database.CheckNewsLetterOnFirestore(country)
	if err != nil {
		fmt.Println(err)
	}

	nws := getStoredNewsLetter(storedNews)

	// check if the interface is null
	if nws.Newsletters == nil || storedNews.Location != country {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		nws, err = getNewsletters(country)
	}

	return nws, err
}

// storeNewsLetter - stores a NewsLetter object to Firestore
func storeNewsLetter(resp structs.NewsLetters, country string) error {
	// populate struct with data to be store
	database.StoredNewsLetter = structs.StoredNewsLetter{
		NewsLetters:  resp,
		Location:     country,
		StoreTime:    time.Now(),
		StoreRefresh: newsLetterDur,
	}
	// save the object to firestore
	err := database.SaveNewsLetterToFirestore(&database.StoredNewsLetter)
	return err
}

// getStoredNewsLetter
func getStoredNewsLetter(storednews structs.StoredNewsLetter) structs.NewsLetters {
	if storednews.NewsLetters.Newsletters == nil {
		return structs.NewsLetters{}
	}
	storednews.StoreRefresh -= time.Since(storednews.StoreTime).Seconds()
	storednews.StoreTime = time.Now()
	database.UpdateTimeFirestore(storednews.FirestoreID, storednews.StoreTime, storednews.StoreRefresh)
	fmt.Println(storednews.StoreRefresh)
	if storednews.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(storednews.FirestoreID)
		return structs.NewsLetters{}
	}
	return storednews.NewsLetters
}

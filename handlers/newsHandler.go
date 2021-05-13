package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"projectGroup23/database"
	"projectGroup23/structs"
	"projectGroup23/utils"

	"google.golang.org/api/iterator"
)

// struct used for the cached data
var storedNewsLetter structs.StoredNewsLetter

var newsLetter structs.NewsLetters

// const for cache duration
const newsLetterDur = 100



// getNewsletters - gets all the newsletters from the api
// this call is only done when no cached data exists at startup
// and when a cached object is deleted after timeout
// TODO - security on API key
// TODO - better error handling
func getNewsletters() structs.NewsLetters {
	fmt.Println("API call made!") // for debugging
	resp, err := http.Get("https://newsapi.org/v2/top-headlines?country=no&apiKey=03b8fc7d5add4ac98eb2330004fbb45c")

	if err != nil {
		fmt.Println(err)
	}
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	jsonRes := string(output)

	newsLetter = utils.PopulateNewsLetters(3, jsonRes)

	// cache the data retrieved from API
	storeNewsLetter(newsLetter)

	return newsLetter
}

// TestEndpoint - just for development, testing that the functionality works correctly
// TODO - remove when not needed anymore
func NewsLetterMainHandler() structs.NewsLetters {
	fmt.Println("NewsletterTest() was run!")
	// use function to retrieve cached newsletter
	nws := getStoredNewsLetter()

	// check if the interface is null
	if nws.Newsletters == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		nws = getNewsletters()
	}

	return nws
}

// storeNewsLetter - stores a NewsLetter object to Firestore
func storeNewsLetter(resp structs.NewsLetters) {
	// populate struct with data to be store
	storedNewsLetter = structs.StoredNewsLetter{
		NewsLetters:   resp,
		StoreTime:    time.Now(),
		StoreRefresh: newsLetterDur,
	}
	// save the object to firestore
	saveNewsLetterToFirestore(&storedNewsLetter)
}

// saveNewsLetterToFirestore - saves an object to firestore
func saveNewsLetterToFirestore(stored *structs.StoredNewsLetter) {
	doc, _, err := database.Client.Collection("cached_resp").Add(database.Ctx, *stored)
	stored.FirestoreID = doc.ID     // storing firestore ID for later use
	fmt.Println(stored.FirestoreID) // confirming the storage of document ID
	if err != nil {
		fmt.Println(err)
	}
}

// GetCachedNewsLetterFromFirestore - global function that runs at startup
// gets all the cached data from firestore
func GetStoredNewsLetterFromFirestore() {
	iter := database.Client.Collection("cached_resp").Documents(database.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		doc.DataTo(&storedNewsLetter)
		storedNewsLetter.FirestoreID = doc.Ref.ID // matching the firestore ID with the one stored
	}
}

// getCachedNewsLetter - used on endpoint to retrieve the cached newsletter
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
func getStoredNewsLetter() structs.NewsLetters {
	if storedNewsLetter.NewsLetters.Newsletters == nil {
		return structs.NewsLetters{}
	}
	storedNewsLetter.StoreRefresh -= time.Since(storedNewsLetter.StoreTime).Seconds()
	storedNewsLetter.StoreTime = time.Now()
	database.UpdateTimeFirestore(storedNewsLetter.FirestoreID, storedNewsLetter.StoreTime, storedNewsLetter.StoreRefresh)
	fmt.Println(storedNewsLetter.StoreRefresh)
	if storedNewsLetter.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(storedNewsLetter.FirestoreID)
		return structs.NewsLetters{}
	}
	return storedNewsLetter.NewsLetters
}

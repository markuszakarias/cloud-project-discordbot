package caching

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"projectGroup23/database"
	"projectGroup23/structs"
	"projectGroup23/utils"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// struct used for the cached data
var c_news CachedNewsLetter

var store_news structs.NewsLetters

var s_newsletter string

// const for cache duration
const c_newsletter_dur = 100

// CachedNewsLetter - struct for a cached newsletter
type CachedNewsLetter struct {
	NewsLetters   structs.NewsLetters
	CachedTime    time.Time
	CachedRefresh float64
	firestoreID   string
}

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
	s_newsletter = string(output)

	store_news = utils.PopulateNewsLetters(3, s_newsletter)

	// cache the data retrieved from API
	cacheNewsLetter(store_news)

	return store_news
}

// TestEndpoint - just for development, testing that the functionality works correctly
// TODO - remove when not needed anymore
func NewsletterTest() structs.NewsLetters {
	fmt.Println("NewsletterTest() was run!")
	// use function to retrieve cached newsletter
	nws := getCachedNewsLetter()

	// check if the interface is null
	if nws.Newsletters == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		nws = getNewsletters()
	}

	return nws
}

// cacheNewsLetter - caches a NewsLetter object to a cache object
// will be stored in firestore
func cacheNewsLetter(resp structs.NewsLetters) {
	// populate struct with data to be cached
	c_news = CachedNewsLetter{
		NewsLetters:   resp,
		CachedTime:    time.Now(),
		CachedRefresh: c_newsletter_dur,
	}
	// save the object on firestore
	saveNewsLetterToFirestore(&c_news)
}

// saveNewsLetterToFirestore - saves an object to firestore
func saveNewsLetterToFirestore(c_save *CachedNewsLetter) {
	doc, _, err := database.Client.Collection("cached_resp").Add(database.Ctx, *c_save)
	c_save.firestoreID = doc.ID     // storing firestore ID for later use
	fmt.Println(c_save.firestoreID) // confirming the storage of document ID
	if err != nil {
		fmt.Println(err)
	}
}

// GetCachedNewsLetterFromFirestore - global function that runs at startup
// gets all the cached data from firestore
func GetCachedNewsLetterFromFirestore() {
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
		c_news.firestoreID = doc.Ref.ID // matching the firestore ID with the one stored
	}
}

// getCachedNewsLetter - used on endpoint to retrieve the cached newsletter
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
func getCachedNewsLetter() structs.NewsLetters {
	if c_news.NewsLetters.Newsletters == nil {
		return structs.NewsLetters{}
	}
	c_news.CachedRefresh -= time.Since(c_news.CachedTime).Seconds()
	c_news.CachedTime = time.Now()
	updateCachedTimeOnNewsLetterFirestore(c_news.firestoreID, c_news.CachedTime, c_news.CachedRefresh)
	fmt.Println(c_news.CachedRefresh)
	if c_news.CachedRefresh <= 0 {
		deleteNewsLetterFromFirestore(c_news.firestoreID)
		return structs.NewsLetters{}
	}
	return c_news.NewsLetters
}

// deleteNewsLetterFromFirestore - deletes an object in firestore based on firestore ID
func deleteNewsLetterFromFirestore(firestoreID string) {
	_, err := database.Client.Collection("cached_resp").Doc(firestoreID).Delete(database.Ctx)
	if err != nil {
		fmt.Println(err)
	}
}

// updateCachedTimeOnNewsLetterFirestore - updates the object in firestore
func updateCachedTimeOnNewsLetterFirestore(firestoreID string, cachedTime time.Time, cachedRefresh float64) {
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

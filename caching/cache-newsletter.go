package caching

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"projectGroup23/firebase"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// struct used for the cached data
var c_news CachedNewsLetter

// struct used to retrieved data from api
var news NewsLetter

// const for cache duration
const c_newsletter_dur = 25

// CachedNewsLetter - struct for a cached newsletter
type CachedNewsLetter struct {
	NewsLetter    NewsLetter
	CachedTime    time.Time
	CachedRefresh float64
	firestoreID   string
}

// Atricle - struct for the data of an article
type Article struct {
	Source      interface{} `json:"source"`
	Author      string      `json:"author"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Url         string      `json:"url"`
	UrlToImage  string      `json:"urlToImage"`
	PublishedAt string      `json:"publishedAt"`
	Content     string      `json:"content"`
}

// NewsLetter - struct to hold the slice of articles
type NewsLetter struct {
	Articles []Article `json:"articles"`
}

// getNewsletters - gets all the newsletters from the api
// this call is only done when no cached data exists at startup
// and when a cached object is deleted after timeout
// TODO - security on API key
// TODO - better error handling
func getNewsletters() interface{} {
	fmt.Println("API call made!") // for debugging
	resp, err := http.Get("https://newsapi.org/v2/top-headlines?country=no&apiKey=03b8fc7d5add4ac98eb2330004fbb45c")

	if err != nil {
		fmt.Println(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&news)
	if err != nil {
		fmt.Println(err)
	}
	// cache the data retrieved from API
	cacheNewsLetter(news)

	// return the populated object
	return news
}

// TestEndpoint - just for development, testing that the functionality works correctly
// TODO - remove when not needed anymore
func NewsletterTest(w http.ResponseWriter, r *http.Request) {
	// use function to retrieve cached newsletter
	nws := getCachedNewsLetter()

	// check if the interface is null
	if nws == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		nws = getNewsletters()
	}
	err := json.NewEncoder(w).Encode(nws)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// cacheNewsLetter - caches a NewsLetter object to a cache object
// will be stored in firestore
func cacheNewsLetter(resp NewsLetter) {
	// populate struct with data to be cached
	c_news = CachedNewsLetter{
		NewsLetter:    resp,
		CachedTime:    time.Now(),
		CachedRefresh: c_newsletter_dur,
	}
	// save the object on firestore
	saveNewsLetterToFirestore(&c_news)
}

// saveNewsLetterToFirestore - saves an object to firestore
func saveNewsLetterToFirestore(c_save *CachedNewsLetter) {
	doc, _, err := firebase.Client.Collection("cached_resp").Add(firebase.Ctx, *c_save)
	c_save.firestoreID = doc.ID     // storing firestore ID for later use
	fmt.Println(c_save.firestoreID) // confirming the storage of document ID
	if err != nil {
		fmt.Println(err)
	}
}

// GetCachedNewsLetterFromFirestore - global function that runs at startup
// gets all the cached data from firestore
func GetCachedNewsLetterFromFirestore() {
	iter := firebase.Client.Collection("cached_resp").Documents(firebase.Ctx)
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
func getCachedNewsLetter() interface{} {
	if c_news.NewsLetter.Articles == nil {
		return nil
	}
	c_news.CachedRefresh -= time.Since(c_news.CachedTime).Seconds()
	c_news.CachedTime = time.Now()
	updateCachedTimeOnNewsLetterFirestore(c_news.firestoreID, c_news.CachedTime, c_news.CachedRefresh)
	fmt.Println(c_news.CachedRefresh)
	if c_news.CachedRefresh <= 0 {
		deleteNewsLetterFromFirestore(c_news.firestoreID)
		return nil
	}
	return c_news.NewsLetter
}

// deleteNewsLetterFromFirestore - deletes an object in firestore based on firestore ID
func deleteNewsLetterFromFirestore(firestoreID string) {
	_, err := firebase.Client.Collection("cached_resp").Doc(firestoreID).Delete(firebase.Ctx)
	if err != nil {
		fmt.Println(err)
	}
}

// updateCachedTimeOnNewsLetterFirestore - updates the object in firestore
func updateCachedTimeOnNewsLetterFirestore(firestoreID string, cachedTime time.Time, cachedRefresh float64) {
	_, err := firebase.Client.Collection("cached_resp").Doc(firestoreID).Update(firebase.Ctx, []firestore.Update{
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

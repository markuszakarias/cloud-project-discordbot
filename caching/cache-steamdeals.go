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
var c_steam CachedSteamdeals

var store_steam structs.Deals

var s_steamdeals string

// const for cache duration
const c_steamdeals_dur = 100

// CachedNewsLetter - struct for a cached newsletter
type CachedSteamdeals struct {
	Steamdeals    structs.Deals
	CachedTime    time.Time
	CachedRefresh float64
	firestoreID   string
}

// getNewsletters - gets all the newsletters from the api
// this call is only done when no cached data exists at startup
// and when a cached object is deleted after timeout
// TODO - security on API key
// TODO - better error handling
func getSteamdeals(command string) structs.Deals {
	fmt.Println("API call made!") // for debugging
	resp, err := http.Get("https://www.cheapshark.com/api/1.0/deals")

	if err != nil {
		fmt.Println(err)
	}
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	s_steamdeals = string(output)

	store_steam = utils.PopulateSteamDeals(s_steamdeals, command, 3)

	// cache the data retrieved from API
	storeSteamdeals(store_steam)

	return store_steam
}

// TestEndpoint - just for development, testing that the functionality works correctly
// TODO - remove when not needed anymore
func SteamdealsTest(command string) structs.Deals {
	// use function to retrieve cached newsletter
	nws := getCachedSteamdeals()

	// check if the interface is null
	if nws.Deals == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		nws = getSteamdeals(command)
	}

	return nws
}

// cacheNewsLetter - caches a NewsLetter object to a cache object
// will be stored in firestore
func storeSteamdeals(resp structs.Deals) {
	// populate struct with data to be cached
	c_steam = CachedSteamdeals{
		Steamdeals:    resp,
		CachedTime:    time.Now(),
		CachedRefresh: c_steamdeals_dur,
	}
	// save the object on firestore
	saveSteamdealsToFirestore(&c_steam)
}

// saveNewsLetterToFirestore - saves an object to firestore
func saveSteamdealsToFirestore(c_save *CachedSteamdeals) {
	doc, _, err := database.Client.Collection("cached_resp").Add(database.Ctx, *c_save)
	c_save.firestoreID = doc.ID     // storing firestore ID for later use
	fmt.Println(c_save.firestoreID) // confirming the storage of document ID
	if err != nil {
		fmt.Println(err)
	}
}

// GetCachedNewsLetterFromFirestore - global function that runs at startup
// gets all the cached data from firestore
func GetSteamdealsFromFirestore() {
	iter := database.Client.Collection("cached_resp").Documents(database.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		doc.DataTo(&c_steam)
		c_steam.firestoreID = doc.Ref.ID // matching the firestore ID with the one stored
	}
}

// getCachedNewsLetter - used on endpoint to retrieve the cached newsletter
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
func getCachedSteamdeals() structs.Deals {
	if c_steam.Steamdeals.Deals == nil {
		return structs.Deals{}
	}
	c_steam.CachedRefresh -= time.Since(c_steam.CachedTime).Seconds()
	c_steam.CachedTime = time.Now()
	updateTimeOnSteamdealsFirestore(c_steam.firestoreID, c_steam.CachedTime, c_steam.CachedRefresh)
	fmt.Println(c_steam.CachedRefresh)
	if c_steam.CachedRefresh <= 0 {
		deleteStoredFromFirestore(c_steam.firestoreID)
		return structs.Deals{}
	}
	return c_steam.Steamdeals
}

// deleteNewsLetterFromFirestore - deletes an object in firestore based on firestore ID
func deleteStoredFromFirestore(firestoreID string) {
	_, err := database.Client.Collection("cached_resp").Doc(firestoreID).Delete(database.Ctx)
	if err != nil {
		fmt.Println(err)
	}
}

// updateCachedTimeOnNewsLetterFirestore - updates the object in firestore
func updateTimeOnSteamdealsFirestore(firestoreID string, cachedTime time.Time, cachedRefresh float64) {
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

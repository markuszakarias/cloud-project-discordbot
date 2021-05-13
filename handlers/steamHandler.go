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
var storedSteamDeals structs.StoredSteamDeals

var steamDeals structs.Deals

// const for cache duration
const steamDealsDur = 100



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
	jsonRes := string(output)

	steamDeals = utils.PopulateSteamDeals(jsonRes, command, 3)

	// cache the data retrieved from API
	storeSteamDeals(steamDeals)

	return steamDeals
}

// TestEndpoint - just for development, testing that the functionality works correctly
// TODO - remove when not needed anymore
func SteamDealsMainHandler(command string) structs.Deals {
	// use function to retrieve cached newsletter
	nws := getStoredSteamDeals()

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
func storeSteamDeals(resp structs.Deals) {
	// populate struct with data to be cached
	storedSteamDeals = structs.StoredSteamDeals {
		SteamDeals:    resp,
		StoreTime:    time.Now(),
		StoreRefresh: steamDealsDur,
	}
	// save the object on firestore
	saveSteamDealsToFirestore(&storedSteamDeals)
}

// saveNewsLetterToFirestore - saves an object to firestore
func saveSteamDealsToFirestore(stored *structs.StoredSteamDeals) {
	doc, _, err := database.Client.Collection("cached_resp").Add(database.Ctx, *stored)
	stored.FirestoreID = doc.ID     // storing firestore ID for later use
	fmt.Println(stored.FirestoreID) // confirming the storage of document ID
	if err != nil {
		fmt.Println(err)
	}
}

// GetCachedNewsLetterFromFirestore - global function that runs at startup
// gets all the cached data from firestore
func GetStoredSteamDealsFromFirestore() {
	iter := database.Client.Collection("cached_resp").Documents(database.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		doc.DataTo(&storedSteamDeals)
		storedSteamDeals.FirestoreID = doc.Ref.ID // matching the firestore ID with the one stored
	}
}

// getCachedNewsLetter - used on endpoint to retrieve the cached newsletter
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
func getStoredSteamDeals() structs.Deals {
	if storedSteamDeals.SteamDeals.Deals == nil {
		return structs.Deals{}
	}
	storedSteamDeals.StoreRefresh -= time.Since(storedSteamDeals.StoreTime).Seconds()
	storedSteamDeals.StoreTime = time.Now()
	database.UpdateTimeFirestore(storedSteamDeals.FirestoreID, storedSteamDeals.StoreTime, storedSteamDeals.StoreRefresh)
	fmt.Println(storedNewsLetter.StoreRefresh)
	if storedSteamDeals.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(storedSteamDeals.FirestoreID)
		return structs.Deals{}
	}
	return storedSteamDeals.SteamDeals
}
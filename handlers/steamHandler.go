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

var steamDeals structs.Deals

// const for cache duration
const steamDealsDur = 100

// getNewsletters - gets all the newsletters from the api
// this call is only done when no cached data exists at startup
// and when a cached object is deleted after timeout
// TODO - security on API key
// TODO - better error handling
func getSteamdeals(command string) (structs.Deals, error) {
	fmt.Println("API call made!") // for debugging
	resp, err := http.Get("https://www.cheapshark.com/api/1.0/deals")

	if err != nil {
		return steamDeals, err
	}
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return steamDeals, err
	}
	jsonRes := string(output)

	steamDeals = utils.PopulateSteamDeals(jsonRes, command, 3)

	// cache the data retrieved from API
	err = storeSteamDeals(steamDeals)

	return steamDeals, err
}

// TestEndpoint - just for development, testing that the functionality works correctly
// TODO - remove when not needed anymore
func SteamDealsMainHandler(command string) (structs.Deals, error) {
	var err error
	// use function to retrieve cached newsletter
	nws := getStoredSteamDeals()

	// check if the interface is null
	if nws.Deals == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		nws, err = getSteamdeals(command)
	}

	return nws, err
}

// cacheNewsLetter - caches a NewsLetter object to a cache object
// will be stored in firestore
func storeSteamDeals(resp structs.Deals) error {
	// populate struct with data to be cached
	database.StoredSteamDeals = structs.StoredSteamDeals{
		SteamDeals:   resp,
		StoreTime:    time.Now(),
		StoreRefresh: steamDealsDur,
	}
	// save the object on firestore
	err := database.SaveSteamDealsToFirestore(&database.StoredSteamDeals)
	return err
}

// saveNewsLetterToFirestore - saves an object to firestore
/*
func saveSteamDealsToFirestore(stored *structs.StoredSteamDeals) error {
	doc, _, err := database.Client.Collection("cached_resp").Add(database.Ctx, *stored)
	stored.FirestoreID = doc.ID     // storing firestore ID for later use
	fmt.Println(stored.FirestoreID) // confirming the storage of document ID
	return err
}
*/

// getCachedNewsLetter - used on endpoint to retrieve the cached newsletter
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
func getStoredSteamDeals() structs.Deals {
	if database.StoredSteamDeals.SteamDeals.Deals == nil {
		return structs.Deals{}
	}
	database.StoredSteamDeals.StoreRefresh -= time.Since(database.StoredSteamDeals.StoreTime).Seconds()
	database.StoredSteamDeals.StoreTime = time.Now()
	database.UpdateTimeFirestore(database.StoredSteamDeals.FirestoreID, database.StoredSteamDeals.StoreTime, database.StoredSteamDeals.StoreRefresh)
	fmt.Println(database.StoredNewsLetter.StoreRefresh)
	if database.StoredSteamDeals.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(database.StoredSteamDeals.FirestoreID)
		return structs.Deals{}
	}
	return database.StoredSteamDeals.SteamDeals
}

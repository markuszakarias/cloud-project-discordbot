package handlers

import (
	"io/ioutil"
	"net/http"
	"time"

	"projectGroup23/database"
	"projectGroup23/structs"
	"projectGroup23/utils"
)

// Struct used to handle data in database management system
var steamDeals structs.Deals

// Const for database storage duration
const steamDealsDur = 100

// getSteamDeals - Requests all Steam deals from the api
//// this call is only done when no stored data exists at startup
//// and when a stored object is deleted after timeout
func getSteamDeals(command string) (structs.Deals, error) {
	resp, err := http.Get("https://www.cheapshark.com/api/1.0/deals")
	if err != nil {
		return steamDeals, err
	}

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return steamDeals, err
	}

	jsonRes := string(output)

	// Populates object with JSON response
	steamDeals = utils.PopulateSteamDeals(jsonRes, command, 5)

	// Store the data retrieved from API
	err = storeSteamDeals(steamDeals)

	return steamDeals, err
}

// SteamDealsMainHandler - Main handler for the !steamdeals command
func SteamDealsMainHandler(command string) (structs.Deals, error) {
	var err error

	// Retrieving possible stored data
	nws := getStoredSteamDeals()

	// Checks if it exists stored data
	if nws.Deals == nil {
		nws, err = getSteamDeals(command)
	}

	return nws, err
}

// storeSteamDeals - Stores a SteamDeals object in the database
func storeSteamDeals(resp structs.Deals) error {
	// Populate struct with data to be stored
	database.StoredSteamDeals = structs.StoredSteamDeals{
		SteamDeals:   resp,
		StoreTime:    time.Now(),
		StoreRefresh: steamDealsDur,
	}
	// Store the object
	err := database.SaveSteamDealsToFirestore(&database.StoredSteamDeals)
	return err
}

// getStoredSteamDeals - Updates timestamps in database storage and retrieves matching object to request
func getStoredSteamDeals() structs.Deals {
	// Checks if it exists a stored response
	if database.StoredSteamDeals.SteamDeals.Deals == nil {
		return structs.Deals{}
	}

	// Calculates timestamp and duration stored in database
	database.StoredSteamDeals.StoreRefresh -= time.Since(database.StoredSteamDeals.StoreTime).Seconds()
	database.StoredSteamDeals.StoreTime = time.Now()

	// Updates new timestamp and duration to Firestore object
	database.UpdateTimeFirestore(database.StoredSteamDeals.FirestoreID, database.StoredSteamDeals.StoreTime, database.StoredSteamDeals.StoreRefresh)

	// If the object storage timer is timed out the object is deleted and then renewed when the next command is called
	if database.StoredSteamDeals.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(database.StoredSteamDeals.FirestoreID)
		return structs.Deals{}
	}
	return database.StoredSteamDeals.SteamDeals
}

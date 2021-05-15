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

// getSteamdeals
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

	steamDeals = utils.PopulateSteamDeals(jsonRes, command, 10)

	// cache the data retrieved from API
	err = storeSteamDeals(steamDeals)

	return steamDeals, err
}

// SteamDealsMainHandler
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

// storeSteamDeals
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

// getStoredSteamDeals
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

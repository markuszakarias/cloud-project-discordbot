package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"projectGroup23/database"
	"projectGroup23/structs"
	"projectGroup23/utils"
	"time"

	"google.golang.org/api/iterator"
)

// struct used for the stored data
var storedMealPlan structs.StoredMealPlan

// struct used to retrieve data from api
var mealPlan structs.MealPlan

// const for cache duration
const mealPlanDur = 50



// getMealPlan - Gets all the meal plans from the api
// this call is only done when no stored data exists at startup
// and when a stored object is deleted after timeout
// TODO - security on API key
// TODO - better error handling
func getMealPlan() structs.MealPlan {
	fmt.Println("API call made!") // for debugging
	resp, err := http.Get("https://api.spoonacular.com/mealplanner/generate?timeFrame=day&apiKey=eeb5e8160efb4bedb1ccc4aa441b0102")

	if err != nil {
		fmt.Println(err)
	}
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	jsonRes := string(output)

	mealPlan = utils.PopulateMealPlan(3, jsonRes)
	// store the data retrieved from API
	storeMealPlan(mealPlan)

	// return the populated object
	return mealPlan
}

// MealPlanMainHandler - Main handler for the !mealplan command
func MealPlanMainHandler() structs.MealPlan {
	// use function to retrieve stored newsletter
	mealPlan = getCachedMealPlanner()

	// check if the interface is null
	if mealPlan.Meals == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		mealPlan = getMealPlan()
	}

	return mealPlan
}

// storeMealPlan - Stores a MealPlan object in the database
func storeMealPlan(resp structs.MealPlan) {
	// populate struct with data to be stored
	storedMealPlan = structs.StoredMealPlan{
		MealPlan:      resp,
		StoreTime:    time.Now(),
		StoreRefresh: mealPlanDur,
	}
	// Store the object to Firestore
	saveMealPlannerToFirestore(&storedMealPlan)
}

// saveNewsLetterToFirestore - saves an object to firestore
func saveMealPlannerToFirestore(stored *structs.StoredMealPlan) {
	doc, _, err := database.Client.Collection("cached_resp").Add(database.Ctx, *stored)
	stored.FirestoreID = doc.ID     // storing firestore ID for later use
	fmt.Println(stored.FirestoreID) // confirming the storage of document ID
	if err != nil {
		fmt.Println(err)
	}
}

// GetCachedNewsLetterFromFirestore - global function that runs at startup
// gets all the cached data from firestore
func GetStoredMealPlannerFromFirestore() {
	iter := database.Client.Collection("cached_resp").Documents(database.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		doc.DataTo(&storedMealPlan)
		storedMealPlan.FirestoreID = doc.Ref.ID // matching the firestore ID with the one stored
	}
}

// getCachedMealPlanner - used on endpoint to retrieve the stored MealPlan
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
func getCachedMealPlanner() structs.MealPlan {
	if storedMealPlan.MealPlan.Meals == nil {
		return structs.MealPlan{}
	}
	storedMealPlan.StoreRefresh -= time.Since(storedMealPlan.StoreTime).Seconds()
	storedMealPlan.StoreTime = time.Now()
	database.UpdateTimeFirestore(storedMealPlan.FirestoreID, storedMealPlan.StoreTime, storedMealPlan.StoreRefresh)
	fmt.Println(storedMealPlan.StoreRefresh)
	if storedMealPlan.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(storedMealPlan.FirestoreID)
		return structs.MealPlan{}
	}
	return storedMealPlan.MealPlan
}





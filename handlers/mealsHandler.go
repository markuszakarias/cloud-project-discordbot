// Package handlers contains all our endpoint handlers. Utility functions used are in package utils.
package handlers

import (
	"io/ioutil"
	"net/http"
	"os"
	"projectGroup23/database"
	"projectGroup23/structs"
	"projectGroup23/utils"
	"time"
)

// Struct used to handle data in database management system
var mealPlan structs.MealPlan

// Const for database storage duration - Every 24 hour
const mealPlanDur = 86_400

// getMealPlan - Requests all meal plans from the api
// this call is only done when no stored data exists at startup
// and when a stored object is deleted after timeout
func getMealPlan() (structs.MealPlan, error) {
	// Get api key from env variable
	apikey := os.Getenv("MEALS_KEY")

	resp, err := http.Get("https://api.spoonacular.com/mealplanner/generate?timeFrame=day&apiKey=" + apikey)
	if err != nil {
		return mealPlan, err
	}

	// Reads response from api request
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return mealPlan, err
	}
	jsonRes := string(output)

	// Populates object with JSON response
	mealPlan = utils.PopulateMealPlan(3, jsonRes)

	// Store the data retrieved from API
	err = storeMealPlan(mealPlan)

	return mealPlan, err
}

// MealPlanMainHandler - Main handler for the !mealplan command
func MealPlanMainHandler() (structs.MealPlan, error) {
	var err error

	// Retrieving possible stored data
	mealPlan = getStoredMealPlanner()

	// Checks if it exists stored data
	if mealPlan.Meals == nil {
		mealPlan, err = getMealPlan()
	}

	return mealPlan, err
}

// storeMealPlan - Stores a MealPlan object in the database
func storeMealPlan(resp structs.MealPlan) error {
	// Populate struct with data to be stored
	database.StoredMealPlan = structs.StoredMealPlan{
		MealPlan:     resp,
		StoreTime:    time.Now(),
		StoreRefresh: mealPlanDur,
	}
	// Store the object
	err := database.SaveMealPlannerToFirestore(&database.StoredMealPlan)
	return err
}

// getStoredMealPlanner - Updates timestamps in database storage and retrieves matching object to request
func getStoredMealPlanner() structs.MealPlan {
	// Checks if it exists a stored response
	if database.StoredMealPlan.MealPlan.Meals == nil {
		return structs.MealPlan{}
	}

	// Calculates timestamp and duration stored in database
	database.StoredMealPlan.StoreRefresh -= time.Since(database.StoredMealPlan.StoreTime).Seconds()
	database.StoredMealPlan.StoreTime = time.Now()

	// Updates new timestamp and duration to Firestore object
	database.UpdateTimeFirestore(database.StoredMealPlan.FirestoreID, database.StoredMealPlan.StoreTime, database.StoredMealPlan.StoreRefresh)

	// If the object storage timer is timed out the object is deleted and then renewed when the next command is called
	if database.StoredMealPlan.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(database.StoredMealPlan.FirestoreID)
	}
	return database.StoredMealPlan.MealPlan
}

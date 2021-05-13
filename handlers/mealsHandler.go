package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"projectGroup23/database"
	"projectGroup23/structs"
	"projectGroup23/utils"
	"time"
)

// struct used to retrieve data from api
var mealPlan structs.MealPlan

// const for cache duration
const mealPlanDur = 100

// getMealPlan - Gets all the meal plans from the api
// this call is only done when no stored data exists at startup
// and when a stored object is deleted after timeout
// TODO - security on API key
// TODO - better error handling
func getMealPlan(apikey string) (structs.MealPlan, error) {
	fmt.Println("API call made!") // for debugging
	resp, err := http.Get("https://api.spoonacular.com/mealplanner/generate?timeFrame=day&apiKey=" + apikey)

	if err != nil {
		return mealPlan, err
	}
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return mealPlan, err
	}
	jsonRes := string(output)

	mealPlan = utils.PopulateMealPlan(3, jsonRes)
	// store the data retrieved from API
	err = storeMealPlan(mealPlan)

	// return the populated object
	return mealPlan, err
}

// MealPlanMainHandler - Main handler for the !mealplan command
func MealPlanMainHandler(apikey string) (structs.MealPlan, error) {
	var err error
	// use function to retrieve stored newsletter
	mealPlan = getCachedMealPlanner()

	// check if the interface is null
	if mealPlan.Meals == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		mealPlan, err = getMealPlan(apikey)
	}

	return mealPlan, err
}

// storeMealPlan - Stores a MealPlan object in the database
func storeMealPlan(resp structs.MealPlan) error {
	// populate struct with data to be stored
	database.StoredMealPlan = structs.StoredMealPlan{
		MealPlan:     resp,
		StoreTime:    time.Now(),
		StoreRefresh: mealPlanDur,
	}
	// Store the object to Firestore
	err := database.SaveMealPlannerToFirestore(&database.StoredMealPlan)
	return err
}

// getCachedMealPlanner - used on endpoint to retrieve the stored MealPlan
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
func getCachedMealPlanner() structs.MealPlan {
	if database.StoredMealPlan.MealPlan.Meals == nil {
		return structs.MealPlan{}
	}
	database.StoredMealPlan.StoreRefresh -= time.Since(database.StoredMealPlan.StoreTime).Seconds()
	database.StoredMealPlan.StoreTime = time.Now()
	database.UpdateTimeFirestore(database.StoredMealPlan.FirestoreID, database.StoredMealPlan.StoreTime, database.StoredMealPlan.StoreRefresh)
	fmt.Println(database.StoredMealPlan.StoreRefresh)
	if database.StoredMealPlan.StoreRefresh <= 0 {
		database.DeleteObjectFromFirestore(database.StoredMealPlan.FirestoreID)
	}
	return database.StoredMealPlan.MealPlan
}

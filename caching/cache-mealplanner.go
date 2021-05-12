package caching

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"projectGroup23/firebase"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// struct used for the cached data
var c_mealplanner CachedMealPlan

// struct used to retrieved data from api
var mealplanner MealPlan

// const for cache duration
const c_mealplanner_dur = 50

type CachedMealPlan struct {
	MealPlan      MealPlan
	CachedTime    time.Time
	CachedRefresh float64
	firestoreID   string
}

// Meal holds information of one meal from the external spoonacular api.
type Meal struct {
	Id             int64  `json:"id"`
	ImageType      string `json:"imageType"`
	Title          string `json:"title"`
	ReadyInMinutes int64  `json:"readyInMinutes"`
	Servings       int    `json:"servings"`
	Url            string `json:"sourceUrl"`
}

// mealPlan is the plan that is given through the discord bot to the user when asked for a meal plan.
type MealPlan struct {
	Meals     []Meal `json:"meals"`
	Nutrients struct {
		Calories      float64 `json:"calories"`
		Protein       float64 `json:"protein"` // should be floating point numbers!
		Fat           float64 `json:"fat"`
		CarboHydrates float64 `json:"carbohydrates"`
	}
}

// getNewsletters - gets all the newsletters from the api
// this call is only done when no cached data exists at startup
// and when a cached object is deleted after timeout
// TODO - security on API key
// TODO - better error handling
func getMealPlanner() interface{} {
	fmt.Println("API call made!") // for debugging
	resp, err := http.Get("https://api.spoonacular.com/mealplanner/generate?timeFrame=day&apiKey=eeb5e8160efb4bedb1ccc4aa441b0102")

	if err != nil {
		fmt.Println(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&mealplanner)
	if err != nil {
		fmt.Println(err)
	}
	// cache the data retrieved from API
	cacheMealPlanner(mealplanner)

	// return the populated object
	return mealplanner
}

// TestEndpoint - just for development, testing that the functionality works correctly
// TODO - remove when not needed anymore
func MealPlannerTest(w http.ResponseWriter, r *http.Request) {
	// use function to retrieve cached newsletter
	mealplan := getCachedMealPlanner()

	// check if the interface is null
	if mealplan == nil {
		fmt.Println("struct is empty")
		// get the newsletters from API if empty
		mealplan = getMealPlanner()
	}
	err := json.NewEncoder(w).Encode(mealplan)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// cacheNewsLetter - caches a NewsLetter object to a cache object
// will be stored in firestore
func cacheMealPlanner(resp MealPlan) {
	// populate struct with data to be cached
	c_mealplanner = CachedMealPlan{
		MealPlan:      resp,
		CachedTime:    time.Now(),
		CachedRefresh: c_mealplanner_dur,
	}
	// save the object on firestore
	saveMealPlannerToFirestore(&c_mealplanner)
}

// saveNewsLetterToFirestore - saves an object to firestore
func saveMealPlannerToFirestore(c_save *CachedMealPlan) {
	doc, _, err := firebase.Client.Collection("cached_resp").Add(firebase.Ctx, *c_save)
	c_save.firestoreID = doc.ID     // storing firestore ID for later use
	fmt.Println(c_save.firestoreID) // confirming the storage of document ID
	if err != nil {
		fmt.Println(err)
	}
}

// GetCachedNewsLetterFromFirestore - global function that runs at startup
// gets all the cached data from firestore
func GetCachedMealPlannerFromFirestore() {
	iter := firebase.Client.Collection("cached_resp").Documents(firebase.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		doc.DataTo(&c_mealplanner)
		c_mealplanner.firestoreID = doc.Ref.ID // matching the firestore ID with the one stored
	}
}

// getCachedNewsLetter - used on endpoint to retrieve the cached newsletter
// will also update the object when timeout has passed
// it also update the fields on the object with data from timeout functionality
func getCachedMealPlanner() interface{} {
	if c_mealplanner.MealPlan.Meals == nil {
		return nil
	}
	c_mealplanner.CachedRefresh -= time.Since(c_mealplanner.CachedTime).Seconds()
	c_mealplanner.CachedTime = time.Now()
	updateCachedTimeOnMealPlannerFirestore(c_mealplanner.firestoreID, c_mealplanner.CachedTime, c_mealplanner.CachedRefresh)
	fmt.Println(c_mealplanner.CachedRefresh)
	if c_mealplanner.CachedRefresh <= 0 {
		deleteMealPlannerFromFirestore(c_mealplanner.firestoreID)
		return nil
	}
	return c_mealplanner.MealPlan
}

// deleteNewsLetterFromFirestore - deletes an object in firestore based on firestore ID
func deleteMealPlannerFromFirestore(firestoreID string) {
	_, err := firebase.Client.Collection("cached_resp").Doc(firestoreID).Delete(firebase.Ctx)
	if err != nil {
		fmt.Println(err)
	}
}

// updateCachedTimeOnNewsLetterFirestore - updates the object in firestore
func updateCachedTimeOnMealPlannerFirestore(firestoreID string, cachedTime time.Time, cachedRefresh float64) {
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

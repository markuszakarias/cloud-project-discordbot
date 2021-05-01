package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"projectGroup23/structs"
	"projectGroup23/utilities"
)

func GetDailyNewsLetter(w http.ResponseWriter, r *http.Request) {
	jsonResponseAsString := utilities.GetNewsApiData(w, r)

	var newsLetter structs.NewsLetters
	articleStruct := utilities.PopulateNewsLetters(newsLetter, jsonResponseAsString)

	data, err := json.Marshal(articleStruct)
	if err != nil {
		log.Printf("%v", "Error during JSON marhsall.")
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetDailyMealPlan(w http.ResponseWriter, r *http.Request) {
	mealPlanResponse := utilities.GetDailyMealPlanData(w, r)

	var mealPlanStruct structs.MealPlan
	mealPlanData := utilities.PopulateMealPlan(mealPlanStruct, mealPlanResponse)

	data, err := json.Marshal(mealPlanData)
	if err != nil {
		log.Printf("%v", "Error during JSON marhsall.")
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

package utilities

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"projectGroup23/structs"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// getNewsApiData returns todays headlines from norwegian media and gives it back in json.
func GetNewsApiData(w http.ResponseWriter, r *http.Request) string {
	url := "https://newsapi.org/v2/top-headlines?country=no&apiKey=cfa7f832f70e41c899bf6b735ef77abf"

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, "bad request to external news api newsapi.org", http.StatusBadRequest)
	}

	r.Header.Add("content-type", "application/json")

	client := &http.Client{}

	// Issue request
	res, err := client.Do(request)
	if err != nil {
		fmt.Errorf("Error in response:", err.Error())
	}

	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("Error when reading response: ", err.Error())
	}

	jsonResponseAsString := string(output)

	return jsonResponseAsString
}

// GetDailyMealPlanData returns the data we use from the food api.
func GetDailyMealPlanData(w http.ResponseWriter, r *http.Request) string {
	url := "https://api.spoonacular.com/mealplanner/generate?timeFrame=day&apiKey=eeb5e8160efb4bedb1ccc4aa441b0102"

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, "bad request to external news api spoonacular.org", http.StatusBadRequest)
	}

	r.Header.Add("content-type", "application/json")

	client := &http.Client{}

	// Issue request
	res, err := client.Do(request)
	if err != nil {
		fmt.Errorf("Error in response:", err.Error())
	}

	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("Error when reading response: ", err.Error())
	}

	jsonResponseAsString := string(output)

	return jsonResponseAsString
}

// populateNewsLetters walks through the response from the newsletter api and creates a
// newsletter json array with 5 newsletters.
func PopulateNewsLetters(paramStruct structs.NewsLetters, jsonResponseString string) structs.NewsLetters {
	for i := 0; i < 5; i++ {
		indexAsString := strconv.Itoa(i) // this counts i as a string from 0-4 throughout the loops iterations.

		authorJson := gjson.Get(jsonResponseString, "articles."+indexAsString+".author")
		publishedAtJson := gjson.Get(jsonResponseString, "articles."+indexAsString+".publishedAt")
		titleJson := gjson.Get(jsonResponseString, "articles."+indexAsString+".title")
		descriptionJson := gjson.Get(jsonResponseString, "articles."+indexAsString+".description")
		urlJson := gjson.Get(jsonResponseString, "articles."+indexAsString+".url")

		author := authorJson.String()
		publishedAt := publishedAtJson.String()
		title := titleJson.String()
		description := descriptionJson.String()
		url := urlJson.String()

		paramStruct.Newsletters[i].Author = author
		paramStruct.Newsletters[i].Date_published = publishedAt
		paramStruct.Newsletters[i].Title = title
		paramStruct.Newsletters[i].Description = description
		paramStruct.Newsletters[i].Url_to_story = url
	}
	return paramStruct
}

// PopulateMealPlan populates a mealPlan struct appropriately.
func PopulateMealPlan(paramStruct structs.MealPlan, jsonResponseString string) structs.MealPlan {
	currentYearMonthDay := time.Now().Format("2006-01-02")
	mealMessage := "Here you go! This is your personal meal plan for today (" + currentYearMonthDay + ")"

	var mealPlanData structs.MealPlan
	mealPlanData.MealMessage = mealMessage

	caloriesJson := gjson.Get(jsonResponseString, "nutrients.calories")
	proteinJson := gjson.Get(jsonResponseString, "nutrients.protein")
	fatJson := gjson.Get(jsonResponseString, "nutrients.fat")
	carboHydratesJson := gjson.Get(jsonResponseString, "nutrients.carbohydrates")

	calories := caloriesJson.Float()
	protein := proteinJson.Float()
	fat := fatJson.Float()
	carboHydrates := carboHydratesJson.Float()

	mealPlanData.Nutrients.Calories = calories
	mealPlanData.Nutrients.Protein = protein
	mealPlanData.Nutrients.Fat = fat
	mealPlanData.Nutrients.CarboHydrates = carboHydrates

	for i := 0; i < 3; i++ {
		indexAsString := strconv.Itoa(i)

		titleJson := gjson.Get(jsonResponseString, "meals."+indexAsString+".title")
		readyInMinutesJson := gjson.Get(jsonResponseString, "meals."+indexAsString+".readyInMinutes")
		urlJson := gjson.Get(jsonResponseString, "meals."+indexAsString+".sourceUrl")

		title := titleJson.String()
		readyInMinutes := readyInMinutesJson.String()
		url := urlJson.String()

		mealPlanData.Meals[i].Title = title
		mealPlanData.Meals[i].ReadyInMinutes = readyInMinutes
		mealPlanData.Meals[i].Url = url
	}
	return mealPlanData
}

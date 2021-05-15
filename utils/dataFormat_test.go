package utils

import (
	"io/ioutil"
	"log"
	"projectGroup23/structs"
	"reflect"
	"testing"
)

func TestMealPopulationFunction(t *testing.T) {

	content, err := ioutil.ReadFile("/projectgroup23/assets/meal.txt")

	if err != nil {
		log.Fatal(err)
	}

	stringOutput := string(content)

	var structWeWant structs.MealPlan

	structWeWant.MealMessage = "Here you go! This is your personal meal plan for today (2021-05-13)"
	structWeWant.Meals = []structs.Meal{
		{Title: "Grape Jelly Breakfast Tarts",
			ReadyInMinutes: "120",
			Url:            "http://www.foodnetwork.com/recipes/food-network-kitchens/grape-jelly-breakfast-tarts-recipe.html"},
		{Title: "Kimchi fried rice",
			ReadyInMinutes: "20",
			Url:            "https://www.bbcgoodfood.com/recipes/kimchi-fried-rice"},
		{Title: "Popcorn & Pretzel Chicken Tenders",
			ReadyInMinutes: "45",
			Url:            "http://www.tasteofhome.com/Recipes/popcorn---pretzel-chicken-tenders"},
	}
	structWeWant.Nutrients.Calories = 1924.26
	structWeWant.Nutrients.CarboHydrates = 269.38
	structWeWant.Nutrients.Fat = 70.32
	structWeWant.Nutrients.Protein = 53.55

	testStruct := PopulateMealPlan(3, stringOutput)

	if reflect.DeepEqual(testStruct, structWeWant) {
		t.Errorf("Population of struct fields was incorrect, got: %+v\n, want: %+v\n", testStruct, structWeWant)
	}
}

func TestNewsletterPopulationFunction(t *testing.T) {

	content, err := ioutil.ReadFile("/projectgroup23/assets/news.txt")
	if err != nil {
		log.Fatal(err)
	}
	stringOutput := string(content)

	var innerStructSliceWeWant structs.NewsLetter
	var structWeWant structs.NewsLetters
	innerStructSliceWeWant.Author = "Håkon Kvam Lyngstad"
	innerStructSliceWeWant.Date_published = "2021-05-13T15:52:37Z"
	innerStructSliceWeWant.Description = "Risør kommune forbyr alle arrangementer og stenger en rekke institusjoner som følge av et coronautbrudd."
	innerStructSliceWeWant.Title = "Innstramming i Risør – flere hundre i karantene – VG - VG"
	innerStructSliceWeWant.UrlToStory = "https://www.vg.no/nyheter/innenriks/i/mBkqPq/innstramming-i-risoer-flere-hundre-i-karantene"

	structWeWant.Newsletters = append(structWeWant.Newsletters, innerStructSliceWeWant)

	testStruct := PopulateNewsLetters(1, stringOutput)

	if reflect.DeepEqual(testStruct, structWeWant) {
		t.Errorf("Population of struct fields was incorrect, got: %+v\n, want: %+v\n", testStruct, structWeWant)
	}
}

func TestSteamPopulationFunction(t *testing.T) {

	content, err := ioutil.ReadFile("/projectgroup23/assets/steam.txt")
	if err != nil {
		log.Fatal(err)
	}
	stringOutput := string(content)

	var innerStructSliceWeWant structs.Deal
	var structWeWant structs.Deals
	innerStructSliceWeWant.Title = "The Lions Song"
	innerStructSliceWeWant.DealID = "Sh%2Fy%2BCAh%2FDfAybb8dWoLnZAYnQVGF3ePV%2F4QM1mOHaQ%3D"
	innerStructSliceWeWant.NormalPrice = "7.99"
	innerStructSliceWeWant.SalePrice = "0.00"
	innerStructSliceWeWant.Savings = "100.000000"

	structWeWant.Deals = append(structWeWant.Deals, innerStructSliceWeWant)

	testStruct := PopulateSteamDeals(stringOutput, "!steamdeals", 1)

	if !reflect.DeepEqual(testStruct, structWeWant) {
		t.Errorf("Population of struct fields was incorrect, got: %+v\n, want: %+v\n", testStruct, structWeWant)
	}
}

func TestWeatherPopulationFunction(t *testing.T) {

	content, err := ioutil.ReadFile("/projectgroup23/assets/weather.txt")
	if err != nil {
		log.Fatal(err)
	}

	stringOutput := string(content)

	var innerStructSliceWeWant structs.WeatherForecast
	var structWeWant structs.WeatherForecasts

	innerStructSliceWeWant.Date = "2021-05-13"
	innerStructSliceWeWant.City = "Oslo"
	innerStructSliceWeWant.Clouds = 75
	innerStructSliceWeWant.Wind = 3.07
	innerStructSliceWeWant.POP = 0.95
	innerStructSliceWeWant.Rain = 7.22
	innerStructSliceWeWant.Snow = 0
	innerStructSliceWeWant.Morning = 9.4
	innerStructSliceWeWant.Night = 15.32
	innerStructSliceWeWant.Day = 20.17
	innerStructSliceWeWant.Main = "Rain"
	innerStructSliceWeWant.Eve = 16.93
	innerStructSliceWeWant.Desc = "moderate rain"

	structWeWant.Forecasts = append(structWeWant.Forecasts, innerStructSliceWeWant)

	testStruct := PopulateWeatherForecast(stringOutput, 1)

	if reflect.DeepEqual(testStruct, structWeWant) {
		t.Errorf("Population of struct fields was incorrect, got: %+v\n, want: %+v\n", testStruct, structWeWant)
	}
}

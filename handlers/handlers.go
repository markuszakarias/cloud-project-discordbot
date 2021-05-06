package handlers

import (
	"projectGroup23/structs"
	"projectGroup23/utilities"
)

func GetDailyNewsLetter() structs.NewsLetters {
	jsonResponseAsString := utilities.GetNewsApiData()
	var newsLetter structs.NewsLetters
	articleStruct := utilities.PopulateNewsLetters(newsLetter, jsonResponseAsString)
	return articleStruct
}

func GetDailyMealPlan() structs.MealPlan {
	mealPlanResponse := utilities.GetDailyMealPlanData()
	var mealPlanStruct structs.MealPlan
	mealPlanData := utilities.PopulateMealPlan(mealPlanStruct, mealPlanResponse)
	return mealPlanData
}

func GetWeatherForecast(days int) structs.WeatherForecasts {
	wfResponse := utilities.GetWeeklyWeatherForecastData(days)
	wf := utilities.PopulateWeatherForecast(wfResponse, days)
	return wf
}

// GetSteamDeals returns deals in the steam store from an external api. It passes along an argument as to how many it should return.
func GetSteamDeals(command string) structs.Deals {
	jsonResponseAsString := utilities.GetSteamDeals()
	dealsData := utilities.GetDeals(jsonResponseAsString, command)
	return dealsData
}

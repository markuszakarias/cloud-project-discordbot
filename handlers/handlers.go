package handlers

import (
	"projectGroup23/structs"
	"projectGroup23/utils"
)

func GetDailyNewsLetter() structs.NewsLetters {
	jsonResponseAsString := utils.GetNewsApiData()
	articleStruct := utils.PopulateNewsLetters(3, jsonResponseAsString)
	return articleStruct
}

func GetDailyMealPlan() structs.MealPlan {
	mealPlanResponse := utils.GetDailyMealPlanData()
	mealPlanData := utils.PopulateMealPlan(3, mealPlanResponse)
	return mealPlanData
}

func GetWeatherForecast(days int) structs.WeatherForecasts {
	wfResponse := utils.GetWeeklyWeatherForecastData(days)
	wf := utils.PopulateWeatherForecast(wfResponse, days)
	return wf
}

// GetSteamDeals returns deals in the steam store from an external api. It passes along an argument as to how many it should return.
func GetSteamDeals(command string) structs.Deals {
	jsonResponseAsString := utils.GetSteamDeals()
	dealsData := utils.PopulateSteamDeals(jsonResponseAsString, command, 3)
	return dealsData
}
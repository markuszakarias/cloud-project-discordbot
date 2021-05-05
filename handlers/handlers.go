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

func GetWeaklyWeatherForecast() structs.WeatherForecast {
	wfResponse := utilities.GetWeeklyWeatherForecastData()
	wf := utilities.PopulateWeatherForecast(wfResponse)
	return wf
}

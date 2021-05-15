// package utils contains functions that populate the structures that contain the reponses from the api.
// It holds the functions that format the various print-outs that happend in the discord client.
// In helper.go there are functions that are intermediary helpers to other functionality in the program.
// Messageformat has all immutable string arrays as function that get used to format helper messages and discord output.
package utils

import (
	"fmt"
	"projectGroup23/structs"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// PopulateNewsLetters walks through the response from the newsletter api and creates a
// newsletter json array with 5 newsletters.
func PopulateNewsLetters(count int, jsonResponseString string) structs.NewsLetters {

	var nws structs.NewsLetter
	var returnNews structs.NewsLetters

	for i := 0; i < 8; i++ {
		indexAsString := strconv.Itoa(i) // this counts i as a string from 0-4 throughout the loops iterations.

		nws.Author = gjson.Get(jsonResponseString, "articles."+indexAsString+".author").String()
		nws.Title = gjson.Get(jsonResponseString, "articles."+indexAsString+".title").String()
		nws.Description = gjson.Get(jsonResponseString, "articles."+indexAsString+".description").String()
		nws.Date_published = gjson.Get(jsonResponseString, "articles."+indexAsString+".publishedAt").String()
		nws.UrlToStory = gjson.Get(jsonResponseString, "articles."+indexAsString+".url").String()

		returnNews.Newsletters = append(returnNews.Newsletters, nws)
	}
	return returnNews
}

// PopulateMealPlan populates a mealPlan struct appropriately.
func PopulateMealPlan(count int, jsonResponseString string) structs.MealPlan {
	currentYearMonthDay := time.Now().Format("2006-01-02")

	var mealPlanData structs.MealPlan
	var meals structs.Meal

	mealPlanData.MealMessage = "Here you go! This is your personal meal plan for today (" + currentYearMonthDay + ")"

	mealPlanData.Nutrients.Calories = gjson.Get(jsonResponseString, "nutrients.calories").Float()
	mealPlanData.Nutrients.Protein = gjson.Get(jsonResponseString, "nutrients.protein").Float()
	mealPlanData.Nutrients.Fat = gjson.Get(jsonResponseString, "nutrients.fat").Float()
	mealPlanData.Nutrients.CarboHydrates = gjson.Get(jsonResponseString, "nutrients.carbohydrates").Float()

	for i := 0; i < count; i++ {
		indexAsString := strconv.Itoa(i)

		meals.Title = gjson.Get(jsonResponseString, "meals."+indexAsString+".title").String()
		meals.ReadyInMinutes = gjson.Get(jsonResponseString, "meals."+indexAsString+".readyInMinutes").String()
		meals.Url = gjson.Get(jsonResponseString, "meals."+indexAsString+".sourceUrl").String()

		mealPlanData.Meals = append(mealPlanData.Meals, meals)
	}
	return mealPlanData
}

// PopulateWeatherForecast populates a WeatherForecasts struct with response from API.
func PopulateWeatherForecast(jsonResponseString string, days int) structs.WeatherForecasts {
	var wf structs.WeatherForecast
	var wfs structs.WeatherForecasts
	cityJson := gjson.Get(jsonResponseString, "city.name")

	for i := 0; i < days; i++ {
		ias := strconv.Itoa(i)

		wf.Date = time.Now().AddDate(0, 0, i).Format("2006-01-02")
		wf.City = cityJson.String()
		wf.Main = gjson.Get(jsonResponseString, "list."+ias+".weather.0.main").String()
		wf.Desc = gjson.Get(jsonResponseString, "list."+ias+".weather.0.description").String()
		wf.Morning = gjson.Get(jsonResponseString, "list."+ias+".temp.morn").Float()
		wf.Day = gjson.Get(jsonResponseString, "list."+ias+".temp.day").Float()
		wf.Eve = gjson.Get(jsonResponseString, "list."+ias+".temp.eve").Float()
		wf.Night = gjson.Get(jsonResponseString, "list."+ias+".temp.night").Float()
		wf.Clouds = gjson.Get(jsonResponseString, "list."+ias+".clouds").Float()
		wf.Wind = gjson.Get(jsonResponseString, "list."+ias+".speed").Float()
		wf.POP = gjson.Get(jsonResponseString, "list."+ias+".pop").Float()
		wf.Rain = gjson.Get(jsonResponseString, "list."+ias+".rain").Float()
		wf.Snow = gjson.Get(jsonResponseString, "list."+ias+".snow").Float()

		wfs.Forecasts = append(wfs.Forecasts, wf)
	}
	return wfs
}

// PopulateSteamDeals fills the deal struct with information about steam deals ready to present with the discord bot.
func PopulateSteamDeals(jsonResponseString string, command string, count int) structs.Deals {
	var deal structs.Deal
	var deals structs.Deals

	for i := 0; i < count; i++ {
		indexAsString := strconv.Itoa(i) // this counts i as a string from 0-4 throughout the loops iterations.

		deal.Title = gjson.Get(jsonResponseString, indexAsString+".title").String()
		deal.DealID = gjson.Get(jsonResponseString, indexAsString+".dealID").String()
		deal.NormalPrice = gjson.Get(jsonResponseString, indexAsString+".normalPrice").String()
		deal.SalePrice = gjson.Get(jsonResponseString, indexAsString+".salePrice").String()
		deal.Savings = gjson.Get(jsonResponseString, indexAsString+".savings").String()
		deals.Deals = append(deals.Deals, deal)
	}
	return deals
}

// WeatherMessageStringFormat utilizes Sprintf to format the discord messages.
// This allows the discord message to be large message instead of multiple smaller messages.
func WeatherMessageStringFormat(stringToPrint []string, day structs.WeatherForecast) string {
	str := fmt.Sprintf(
		"%s%s\n %s%s\n %s%s%s\n %s%s%s\n %s%s%s\n %s%s\n %s%s%s\n %s%s%s\n %s\n %s%s%s\n %s%s%s\n %s%s%s\n %s%s%s\n",
		stringToPrint[0], day.Date, stringToPrint[1], day.City, day.Main,
		stringToPrint[2], day.Desc, stringToPrint[3], fmt.Sprint(day.Clouds),
		stringToPrint[15], stringToPrint[4], fmt.Sprint(day.Wind), stringToPrint[5],
		stringToPrint[7], fmt.Sprint(day.POP), stringToPrint[8], fmt.Sprint(day.Rain), stringToPrint[5],
		stringToPrint[9], fmt.Sprint(day.Snow), stringToPrint[5], stringToPrint[10],
		stringToPrint[11], fmt.Sprint(day.Morning), stringToPrint[6],
		stringToPrint[12], fmt.Sprint(day.Day), stringToPrint[6],
		stringToPrint[13], fmt.Sprint(day.Eve), stringToPrint[6],
		stringToPrint[14], fmt.Sprint(day.Night), stringToPrint[6])
	return str
}

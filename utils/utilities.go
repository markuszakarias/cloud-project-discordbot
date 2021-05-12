package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"projectGroup23/structs"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// getNewsApiData returns todays headlines from norwegian media and gives it back in json.
func GetNewsApiData() string {
	url := "https://newsapi.org/v2/top-headlines?country=no&apiKey=cfa7f832f70e41c899bf6b735ef77abf"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("Error in response:", err.Error())
	}

	defer resp.Body.Close()

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error when reading response: ", err.Error())
	}

	jsonResponseAsString := string(output)

	fmt.Println(jsonResponseAsString)

	return jsonResponseAsString
}

// GetSteamDeals returns current steam deals and gives it back in json.
func GetSteamDeals() string {
	url := "https://www.cheapshark.com/api/1.0/deals"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("Error in response:", err.Error())
	}

	defer resp.Body.Close()

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error when reading response: ", err.Error())
	}

	jsonResponseAsString := string(output)

	return jsonResponseAsString
}

// GetDailyMealPlanData returns the data we use from the food api.
func GetDailyMealPlanData() string {
	url := "https://api.spoonacular.com/mealplanner/generate?timeFrame=day&apiKey=eeb5e8160efb4bedb1ccc4aa441b0102"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("Error in response:", err.Error())
	}

	defer resp.Body.Close()

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error when reading response: ", err.Error())
	}

	jsonResponseAsString := string(output)

	return jsonResponseAsString
}

func GetIPLocation() structs.IPLocation {
	url := "https://ipwhois.app/json/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("Error in response: ", err.Error())
	}

	defer resp.Body.Close()

	var ipp structs.IPLocation
	err = json.NewDecoder(resp.Body).Decode(&ipp)
	if err != nil {
		fmt.Errorf("Error in JSON decoding: ", err.Error())
	}

	return ipp
}

func GetWeeklyWeatherForecastData(days int) string {
	ipp := GetIPLocation()
	units := "metric"
	cnt := strconv.Itoa(days)
	apikey := "f6a8e67b1a5f1d5be2bffe4d461cc155" //TODO - Secure API key

	url := "https://api.openweathermap.org/data/2.5/forecast/daily?q=" + ipp.City +
		"&units=" + units + "&cnt=" + cnt + "&appid=" + apikey

	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("Error in response: ", err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error when reading response: ", err.Error())
	}

	return string(body)
}

// populateNewsLetters walks through the response from the newsletter api and creates a
// newsletter json array with 5 newsletters.
func PopulateNewsLetters(count int, jsonResponseString string) structs.NewsLetters {

	var nws structs.NewsLetter
	var returnNews structs.NewsLetters

	for i := 0; i < count; i++ {
		indexAsString := strconv.Itoa(i) // this counts i as a string from 0-4 throughout the loops iterations.

		nws.Author = gjson.Get(jsonResponseString, "articles."+indexAsString+".author").String()
		nws.Title = gjson.Get(jsonResponseString, "articles."+indexAsString+".title").String()
		nws.Description = gjson.Get(jsonResponseString, "articles."+indexAsString+".description").String()
		nws.Date_published = gjson.Get(jsonResponseString, "articles."+indexAsString+".publishedAt").String()
		nws.Url_to_story = gjson.Get(jsonResponseString, "articles."+indexAsString+".url").String()

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

// GetDeals fills the deal struct with information about steam deals ready to present with the discord bot.
func PopulateSteamDeals(jsonResponseString string, command string, count int) structs.Deals {
	var deal structs.Deal
	var deals structs.Deals

	if len(command) == 11 {
		for i := 0; i < count; i++ {
			indexAsString := strconv.Itoa(i) // this counts i as a string from 0-4 throughout the loops iterations.

			deal.Title = gjson.Get(jsonResponseString, indexAsString+".title").String()
			deal.DealID = gjson.Get(jsonResponseString, indexAsString+".dealID").String()
			deal.NormalPrice = gjson.Get(jsonResponseString, indexAsString+".normalPrice").String()
			deal.SalePrice = gjson.Get(jsonResponseString, indexAsString+".salePrice").String()
			deal.Savings = gjson.Get(jsonResponseString, indexAsString+".savings").String()
			deal.MetacriticScore = gjson.Get(jsonResponseString, indexAsString+".metacriticScore").String()
			deal.SteamRatingText = gjson.Get(jsonResponseString, indexAsString+".steamRatingText").String()
			deal.SteamRatingPercent = gjson.Get(jsonResponseString, indexAsString+".steamRatingPercent").String()
			deal.SteamRatingCount = gjson.Get(jsonResponseString, indexAsString+".steamRatingCount").String()

			if deal.MetacriticScore == "0" {
				deal.MetacriticScore = "could not fetch metacritic score"
			}
			if deal.SteamRatingCount == "0" {
				deal.SteamRatingCount = "could not grab SteamRatingCount"
			}
			if deal.SteamRatingPercent == "0" {
				deal.SteamRatingPercent = "could not fetch SteamRatingPercent"
			}
			if len(deal.SteamRatingText) == 0 {
				deal.SteamRatingText = "could not fetch SteamRatingText"
			}

			deals.Deals = append(deals.Deals, deal)
		}
	} else {
		strArr := []rune(command)
		amountOfDealsToGet := string(strArr[12])
		numberOfIteration, _ := strconv.Atoi(amountOfDealsToGet)
		for i := 0; i < numberOfIteration; i++ {
			indexAsString := strconv.Itoa(i) // this counts i as a string from 0-4 throughout the loops iterations.

			deal.Title = gjson.Get(jsonResponseString, indexAsString+".title").String()
			deal.DealID = gjson.Get(jsonResponseString, indexAsString+".dealID").String()
			deal.NormalPrice = gjson.Get(jsonResponseString, indexAsString+".normalPrice").String()
			deal.SalePrice = gjson.Get(jsonResponseString, indexAsString+".salePrice").String()
			deal.Savings = gjson.Get(jsonResponseString, indexAsString+".savings").String()
			deal.MetacriticScore = gjson.Get(jsonResponseString, indexAsString+".metacriticScore").String()
			deal.SteamRatingText = gjson.Get(jsonResponseString, indexAsString+".steamRatingText").String()
			deal.SteamRatingPercent = gjson.Get(jsonResponseString, indexAsString+".steamRatingPercent").String()
			deal.SteamRatingCount = gjson.Get(jsonResponseString, indexAsString+".steamRatingCount").String()

			if deal.MetacriticScore == "0" {
				deal.MetacriticScore = "Unavailable to fetch metacritic score"
			}
			if deal.SteamRatingCount == "0" {
				deal.SteamRatingCount = "Unavailable to grab SteamRatingCount"
			}
			if deal.SteamRatingPercent == "0" {
				deal.SteamRatingPercent = "Unavailable to fetch SteamRatingPercent"
			}
			if deal.SteamRatingText == "null" {
				deal.SteamRatingText = "Unavailable to fetch SteamRatingText"
			}
			deals.Deals = append(deals.Deals, deal)
		}
	}
	return deals
}

func CheckIfSameDate(date, date2 time.Time) bool {
	y, m, d := date.Date()
	y2, m2, d2 := date2.Date()
	return y == y2 && m == m2 && d == d2
}

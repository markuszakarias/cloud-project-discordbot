package utilities

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

func PopulateWeatherForecast(jsonResponseString string, days int) structs.WeatherForecasts {
	var wf structs.WeatherForecast
	var wfs structs.WeatherForecasts
	cityJson := gjson.Get(jsonResponseString, "city.name")

	for i := 0; i < days; i++ {
		ias := strconv.Itoa(i)

		date := time.Now().AddDate(0, 0, i).Format("2006-01-02")
		mainJson := gjson.Get(jsonResponseString, "list."+ias+".weather.0.main")
		descJson := gjson.Get(jsonResponseString, "list."+ias+".weather.0.description")
		mornJson := gjson.Get(jsonResponseString, "list."+ias+".temp.morn")
		dayJson := gjson.Get(jsonResponseString, "list."+ias+".temp.day")
		eveJson := gjson.Get(jsonResponseString, "list."+ias+".temp.eve")
		nightJson := gjson.Get(jsonResponseString, "list."+ias+".temp.night")
		cloudsJson := gjson.Get(jsonResponseString, "list."+ias+".clouds")
		windJson := gjson.Get(jsonResponseString, "list."+ias+".speed")
		popJson := gjson.Get(jsonResponseString, "list."+ias+".pop")
		rainJson := gjson.Get(jsonResponseString, "list."+ias+".rain")
		snowJson := gjson.Get(jsonResponseString, "list."+ias+".snow")

		wf.Date = date
		wf.City = cityJson.String()
		wf.Main = mainJson.String()
		wf.Desc = descJson.String()
		wf.Morning = mornJson.Float()
		wf.Day = dayJson.Float()
		wf.Eve = eveJson.Float()
		wf.Night = nightJson.Float()
		wf.Clouds = cloudsJson.Float()
		wf.Wind = windJson.Float()
		wf.POP = popJson.Float()
		wf.Rain = rainJson.Float()
		wf.Snow = snowJson.Float()

		wfs.Forecasts = append(wfs.Forecasts, wf)
	}
	return wfs
}

// GetDeals fills the deal struct with information about steam deals ready to present with the discord bot.
func GetDeals(jsonResponseString string, command string) structs.Deals {
	var deal structs.Deal
	var deals structs.Deals

	if len(command) == 11 {
		for i := 0; i < 3; i++ {
			indexAsString := strconv.Itoa(i) // this counts i as a string from 0-4 throughout the loops iterations.

			titleJson := gjson.Get(jsonResponseString, indexAsString+".title")
			dealIDJson := gjson.Get(jsonResponseString, indexAsString+".dealID")
			normalPriceJson := gjson.Get(jsonResponseString, indexAsString+".normalPrice")
			salePriceJson := gjson.Get(jsonResponseString, indexAsString+".salePrice")
			savingsJson := gjson.Get(jsonResponseString, indexAsString+".savings")
			MetacriticScoreJson := gjson.Get(jsonResponseString, indexAsString+".metacriticScore")
			SteamRatingTextJson := gjson.Get(jsonResponseString, indexAsString+".steamRatingText")
			SteamRatingPercentJson := gjson.Get(jsonResponseString, indexAsString+".steamRatingPercent")
			SteamRatingCountJson := gjson.Get(jsonResponseString, indexAsString+".steamRatingCount")

			title := titleJson.String()
			dealID := dealIDJson.String()
			normalPrice := normalPriceJson.String()
			salePrice := salePriceJson.String()
			savings := savingsJson.String()
			metacriticScore := MetacriticScoreJson.String()
			steamRatingText := SteamRatingTextJson.String()
			steamRatingPercent := SteamRatingPercentJson.String()
			steamRatingCount := SteamRatingCountJson.String()

			deal.Title = title
			deal.DealID = dealID
			deal.NormalPrice = normalPrice
			deal.SalePrice = salePrice
			deal.Savings = savings
			deal.MetacriticScore = metacriticScore
			deal.SteamRatingText = steamRatingText
			deal.SteamRatingPercent = steamRatingPercent
			deal.SteamRatingCount = steamRatingCount

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

			titleJson := gjson.Get(jsonResponseString, indexAsString+".title")
			dealIDJson := gjson.Get(jsonResponseString, indexAsString+".dealID")
			normalPriceJson := gjson.Get(jsonResponseString, indexAsString+".normalPrice")
			salePriceJson := gjson.Get(jsonResponseString, indexAsString+".salePrice")
			savingsJson := gjson.Get(jsonResponseString, indexAsString+".savings")
			MetacriticScoreJson := gjson.Get(jsonResponseString, indexAsString+".metacriticScore")
			SteamRatingTextJson := gjson.Get(jsonResponseString, indexAsString+".steamRatingText")
			SteamRatingPercentJson := gjson.Get(jsonResponseString, indexAsString+".steamRatingPercent")
			SteamRatingCountJson := gjson.Get(jsonResponseString, indexAsString+".steamRatingCount")

			title := titleJson.String()
			dealID := dealIDJson.String()
			normalPrice := normalPriceJson.String()
			salePrice := salePriceJson.String()
			savings := savingsJson.String()
			metacriticScore := MetacriticScoreJson.String()
			steamRatingText := SteamRatingTextJson.String()
			steamRatingPercent := SteamRatingPercentJson.String()
			steamRatingCount := SteamRatingCountJson.String()

			deal.Title = title
			deal.DealID = dealID
			deal.NormalPrice = normalPrice
			deal.SalePrice = salePrice
			deal.Savings = savings
			deal.MetacriticScore = metacriticScore
			deal.SteamRatingText = steamRatingText
			deal.SteamRatingPercent = steamRatingPercent
			deal.SteamRatingCount = steamRatingCount

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
